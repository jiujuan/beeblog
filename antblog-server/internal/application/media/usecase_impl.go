package media

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"

	"go.uber.org/fx"
	"go.uber.org/zap"
	_ "golang.org/x/image/webp"

	domain "antblog/internal/domain/media"
	apperrors "antblog/pkg/errors"
	"antblog/pkg/config"
	"antblog/pkg/utils"
)

// ─── 依赖 ────────────────────────────────────────────────────────────────────

type mediaUseCase struct {
	repo         domain.IMediaRepository
	storage      domain.IStorageDriver
	svc          domain.IDomainService
	logger       *zap.Logger
	allowedTypes []string
	maxSize      int64
}

// NewMediaUseCase 创建媒体用例（fx provider）
// 直接注入 *config.Config，避免 fx 命名值注入的复杂性
func NewMediaUseCase(
	repo domain.IMediaRepository,
	storage domain.IStorageDriver,
	svc domain.IDomainService,
	cfg *config.Config,
	logger *zap.Logger,
) IMediaUseCase {
	return &mediaUseCase{
		repo:         repo,
		storage:      storage,
		svc:          svc,
		logger:       logger,
		allowedTypes: cfg.Upload.AllowedTypes,
		maxSize:      cfg.Upload.MaxSize,
	}
}

// ─── 上传 ────────────────────────────────────────────────────────────────────

func (uc *mediaUseCase) Upload(
	ctx context.Context,
	uploaderID uint64,
	fileHeader *multipart.FileHeader,
	articleID *uint64,
) (*MediaResp, error) {

	// 1. 读取文件内容
	f, err := fileHeader.Open()
	if err != nil {
		return nil, apperrors.New(apperrors.CodeMediaUploadFailed, "无法打开上传文件")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, apperrors.New(apperrors.CodeMediaUploadFailed, "读取文件内容失败")
	}

	originalName := fileHeader.Filename

	// 2. 领域校验（大小、MIME 类型白名单）
	if err = uc.svc.ValidateUpload(originalName, data, uc.allowedTypes, uc.maxSize); err != nil {
		return nil, err
	}

	mimeType := uc.svc.DetectMimeType(data, originalName)

	// 3. SHA256 去重：同一文件内容已存在则直接返回已有记录
	hash := domain.Sha256Hex(data)
	if existing, findErr := uc.repo.FindByHash(ctx, hash); findErr == nil && existing != nil {
		uc.logger.Info("media dedup hit",
			zap.String("hash", hash), zap.Uint64("existing_id", existing.ID))
		// 若请求指定绑定文章，同步更新绑定关系
		if articleID != nil && (existing.ArticleID == nil || *existing.ArticleID != *articleID) {
			existing.BindArticle(*articleID)
			_ = uc.repo.Update(ctx, existing)
		}
		return toMediaResp(existing), nil
	}

	// 4. 读取图片尺寸（仅图片类型，失败不中断）
	width, height := 0, 0
	if isImageMime(mimeType) {
		width, height = readImageDimensions(data)
	}

	// 5. 调用存储驱动写文件
	storagePath, url, err := uc.storage.Save(originalName, data)
	if err != nil {
		uc.logger.Error("storage save failed", zap.String("name", originalName), zap.Error(err))
		return nil, apperrors.New(apperrors.CodeMediaUploadFailed, "文件存储失败")
	}

	// 6. 构建领域实体
	m, err := uc.svc.BuildMedia(uploaderID, originalName, data, storagePath, url, mimeType, width, height)
	if err != nil {
		_ = uc.storage.Delete(storagePath) // 回滚：删除已写入文件
		return nil, err
	}
	m.ArticleID = articleID

	// 7. 持久化到数据库
	created, err := uc.repo.Create(ctx, m)
	if err != nil {
		_ = uc.storage.Delete(storagePath) // 回滚
		uc.logger.Error("media create failed", zap.Error(err))
		return nil, apperrors.ErrInternalError(err)
	}

	uc.logger.Info("media uploaded",
		zap.Uint64("id", created.ID),
		zap.String("name", originalName),
		zap.String("mime", mimeType),
		zap.Int64("size", created.FileSize),
	)
	return toMediaResp(created), nil
}

// ─── 查询 ────────────────────────────────────────────────────────────────────

func (uc *mediaUseCase) GetMedia(ctx context.Context, id uint64) (*MediaResp, error) {
	m, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrMediaNotFound()
	}
	return toMediaResp(m), nil
}

func (uc *mediaUseCase) ListMyMedia(
	ctx context.Context, uploaderID uint64, page, pageSize int,
) ([]*MediaResp, int64, error) {
	page = utils.NormalizePage(page)
	pageSize = utils.NormalizePageSize(pageSize)
	list, total, err := uc.repo.FindByUploaderID(ctx, uploaderID, page, pageSize)
	if err != nil {
		return nil, 0, apperrors.ErrInternalError(err)
	}
	return toMediaRespList(list), total, nil
}

func (uc *mediaUseCase) AdminListMedia(
	ctx context.Context, req *AdminListMediaReq,
) ([]*MediaResp, int64, error) {
	req.Page = utils.NormalizePage(req.Page)
	req.PageSize = utils.NormalizePageSize(req.PageSize)

	filter := &domain.AdminFilter{
		MimeType: req.MimeType,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if req.UploaderID > 0 {
		filter.UploaderID = &req.UploaderID
	}
	if req.ArticleID > 0 {
		filter.ArticleID = &req.ArticleID
	}

	list, total, err := uc.repo.AdminFind(ctx, filter)
	if err != nil {
		return nil, 0, apperrors.ErrInternalError(err)
	}
	return toMediaRespList(list), total, nil
}

// ─── 绑定文章 ────────────────────────────────────────────────────────────────

func (uc *mediaUseCase) BindArticle(
	ctx context.Context, mediaID uint64, req *BindArticleReq,
) (*MediaResp, error) {
	m, err := uc.repo.FindByID(ctx, mediaID)
	if err != nil {
		return nil, apperrors.ErrMediaNotFound()
	}

	if req.ArticleID == nil {
		m.UnbindArticle()
	} else {
		m.BindArticle(*req.ArticleID)
	}

	if err = uc.repo.Update(ctx, m); err != nil {
		return nil, apperrors.ErrInternalError(err)
	}
	return toMediaResp(m), nil
}

// ─── 删除 ────────────────────────────────────────────────────────────────────

func (uc *mediaUseCase) DeleteMedia(ctx context.Context, id uint64) error {
	m, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.ErrMediaNotFound()
	}

	if err = uc.repo.Delete(ctx, id); err != nil {
		return apperrors.ErrInternalError(err)
	}

	// 物理文件删除允许失败（不回滚数据库记录，可后台补偿）
	if delErr := uc.storage.Delete(m.StoragePath); delErr != nil {
		uc.logger.Warn("delete physical file failed",
			zap.String("path", m.StoragePath), zap.Error(delErr))
	}

	uc.logger.Info("media deleted", zap.Uint64("id", id), zap.String("path", m.StoragePath))
	return nil
}

// ─── 私有辅助 ────────────────────────────────────────────────────────────────

func readImageDimensions(data []byte) (width, height int) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}

func isImageMime(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return true
	}
	return false
}

func humanFileSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.1f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// ─── 映射 ────────────────────────────────────────────────────────────────────

func toMediaResp(m *domain.Media) *MediaResp {
	return &MediaResp{
		ID:            m.ID,
		UploaderID:    m.UploaderID,
		ArticleID:     m.ArticleID,
		OriginalName:  m.OriginalName,
		URL:           m.URL,
		MimeType:      m.MimeType,
		FileSize:      m.FileSize,
		FileSizeHuman: humanFileSize(m.FileSize),
		Width:         m.Width,
		Height:        m.Height,
		Hash:          m.Hash,
		CreatedAt:     m.CreatedAt,
	}
}

func toMediaRespList(list []*domain.Media) []*MediaResp {
	resp := make([]*MediaResp, 0, len(list))
	for _, m := range list {
		resp = append(resp, toMediaResp(m))
	}
	return resp
}

// ─── fx Module ───────────────────────────────────────────────────────────────

// Module fx 媒体应用模块
var Module = fx.Options(
	fx.Provide(
		NewMediaUseCase,
		domain.NewDomainService,
	),
)
