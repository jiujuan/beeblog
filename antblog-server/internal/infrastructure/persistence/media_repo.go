package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	domain "antblog/internal/domain/media"
	"antblog/internal/infrastructure/persistence/model"
	apperrors "antblog/pkg/errors"
)

// mediaRepository IMediaRepository 的 GORM 实现
type mediaRepository struct {
	db *gorm.DB
}

// NewMediaRepository 创建媒体仓储
func NewMediaRepository(db *gorm.DB) domain.IMediaRepository {
	return &mediaRepository{db: db}
}

// ─── 单条操作 ────────────────────────────────────────────────────────────────

func (r *mediaRepository) Create(ctx context.Context, m *domain.Media) (*domain.Media, error) {
	model := mediaDomainToModel(m)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return mediaModelToDomain(model), nil
}

func (r *mediaRepository) FindByID(ctx context.Context, id uint64) (*domain.Media, error) {
	var m model.Media
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrMediaNotFound()
		}
		return nil, err
	}
	return mediaModelToDomain(&m), nil
}

func (r *mediaRepository) FindByHash(ctx context.Context, hash string) (*domain.Media, error) {
	var m model.Media
	if err := r.db.WithContext(ctx).Where("hash = ?", hash).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrMediaNotFound()
		}
		return nil, err
	}
	return mediaModelToDomain(&m), nil
}

func (r *mediaRepository) Update(ctx context.Context, m *domain.Media) error {
	return r.db.WithContext(ctx).
		Model(&model.Media{}).
		Where("id = ?", m.ID).
		Updates(map[string]any{
			"article_id": m.ArticleID,
		}).Error
}

func (r *mediaRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Media{}, "id = ?", id).Error
}

// ─── 列表查询 ────────────────────────────────────────────────────────────────

func (r *mediaRepository) FindByUploaderID(
	ctx context.Context, uploaderID uint64, page, pageSize int,
) ([]*domain.Media, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&model.Media{}).Where("uploader_id = ?", uploaderID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	var list []model.Media
	offset := (page - 1) * pageSize
	if err := q.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return mediaModelsToDomainsSlice(list), total, nil
}

func (r *mediaRepository) FindByArticleID(
	ctx context.Context, articleID uint64,
) ([]*domain.Media, error) {
	var list []model.Media
	if err := r.db.WithContext(ctx).
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return mediaModelsToDomainsSlice(list), nil
}

func (r *mediaRepository) AdminFind(
	ctx context.Context, filter *domain.AdminFilter,
) ([]*domain.Media, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Media{})

	if filter.UploaderID != nil {
		q = q.Where("uploader_id = ?", *filter.UploaderID)
	}
	if filter.ArticleID != nil {
		q = q.Where("article_id = ?", *filter.ArticleID)
	}
	if filter.MimeType != "" {
		q = q.Where("mime_type = ?", filter.MimeType)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	var list []model.Media
	offset := (filter.Page - 1) * filter.PageSize
	if err := q.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return mediaModelsToDomainsSlice(list), total, nil
}

// ─── 批量操作 ────────────────────────────────────────────────────────────────

func (r *mediaRepository) DeleteByUploaderID(ctx context.Context, uploaderID uint64) error {
	return r.db.WithContext(ctx).
		Where("uploader_id = ?", uploaderID).
		Delete(&model.Media{}).Error
}

// ─── 模型映射 ────────────────────────────────────────────────────────────────

func mediaDomainToModel(m *domain.Media) *model.Media {
	return &model.Media{
		ID:           m.ID,
		UploaderID:   m.UploaderID,
		ArticleID:    m.ArticleID,
		OriginalName: m.OriginalName,
		StoragePath:  m.StoragePath,
		URL:          m.URL,
		MimeType:     m.MimeType,
		FileSize:     m.FileSize,
		Width:        m.Width,
		Height:       m.Height,
		Hash:         m.Hash,
	}
}

func mediaModelToDomain(m *model.Media) *domain.Media {
	return &domain.Media{
		ID:           m.ID,
		UploaderID:   m.UploaderID,
		ArticleID:    m.ArticleID,
		OriginalName: m.OriginalName,
		StoragePath:  m.StoragePath,
		URL:          m.URL,
		MimeType:     m.MimeType,
		FileSize:     m.FileSize,
		Width:        m.Width,
		Height:       m.Height,
		Hash:         m.Hash,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func mediaModelsToDomainsSlice(list []model.Media) []*domain.Media {
	result := make([]*domain.Media, 0, len(list))
	for i := range list {
		result = append(result, mediaModelToDomain(&list[i]))
	}
	return result
}
