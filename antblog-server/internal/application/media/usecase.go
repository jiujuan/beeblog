package media

import (
	"context"
	"mime/multipart"
)

// IMediaUseCase 媒体用例接口
type IMediaUseCase interface {
	// ── 后台管理接口（需 Admin 权限）──────────────────────────────────────────

	// Upload 上传媒体文件（校验、去重、存储、持久化）
	// fileHeader: gin 从 multipart 解析的文件头
	Upload(ctx context.Context, uploaderID uint64, fileHeader *multipart.FileHeader, articleID *uint64) (*MediaResp, error)

	// GetMedia 按 ID 获取媒体详情
	GetMedia(ctx context.Context, id uint64) (*MediaResp, error)

	// ListMyMedia 查询当前登录用户上传的媒体（分页）
	ListMyMedia(ctx context.Context, uploaderID uint64, page, pageSize int) ([]*MediaResp, int64, error)

	// AdminListMedia 后台媒体列表（多条件过滤）
	AdminListMedia(ctx context.Context, req *AdminListMediaReq) ([]*MediaResp, int64, error)

	// BindArticle 将媒体绑定到文章（或解绑：articleID 传 nil）
	BindArticle(ctx context.Context, mediaID uint64, req *BindArticleReq) (*MediaResp, error)

	// DeleteMedia 删除媒体资源（软删除记录 + 物理文件删除）
	DeleteMedia(ctx context.Context, id uint64) error
}
