package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	domain "antblog/internal/domain/category"
	"antblog/internal/infrastructure/persistence/model"
	apperrors "antblog/pkg/errors"
)

// categoryRepository ICategoryRepository 的 GORM 实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) domain.ICategoryRepository {
	return &categoryRepository{db: db}
}

// Create 创建分类
func (r *categoryRepository) Create(ctx context.Context, c *domain.Category) (*domain.Category, error) {
	m := categoryDomainToModel(c)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return categoryModelToDomain(m), nil
}

// FindByID 按 ID 查询分类
func (r *categoryRepository) FindByID(ctx context.Context, id uint64) (*domain.Category, error) {
	var m model.Category
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCategoryNotFound()
		}
		return nil, err
	}
	return categoryModelToDomain(&m), nil
}

// FindBySlug 按 Slug 查询分类
func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var m model.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrCategoryNotFound()
		}
		return nil, err
	}
	return categoryModelToDomain(&m), nil
}

// FindAll 查询所有分类，按 sort_order 降序
func (r *categoryRepository) FindAll(ctx context.Context) ([]*domain.Category, error) {
	var list []model.Category
	if err := r.db.WithContext(ctx).
		Model(&model.Category{}).
		Select(`
			categories.id,
			categories.name,
			categories.slug,
			categories.description,
			categories.cover,
			categories.sort_order,
			categories.created_at,
			categories.updated_at,
			(
				SELECT COUNT(1)
				FROM articles
				WHERE articles.category_id = categories.id
				  AND articles.status = ?
				  AND articles.deleted_at IS NULL
			) AS article_count
		`, 2).
		Order("sort_order DESC, id ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.Category, 0, len(list))
	for i := range list {
		result = append(result, categoryModelToDomain(&list[i]))
	}
	return result, nil
}

// Update 更新分类
func (r *categoryRepository) Update(ctx context.Context, c *domain.Category) error {
	return r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", c.ID).
		Updates(map[string]any{
			"name":        c.Name,
			"slug":        c.Slug,
			"description": c.Description,
			"cover":       c.Cover,
			"sort_order":  c.SortOrder,
		}).Error
}

// Delete 软删除分类
func (r *categoryRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, "id = ?", id).Error
}

// ExistsBySlug 检查 Slug 是否存在（excludeID=0 表示不排除任何记录）
func (r *categoryRepository) ExistsBySlug(ctx context.Context, slug string, excludeID uint64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.Category{}).Where("slug = ?", slug)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// ExistsByName 检查名称是否存在（excludeID=0 表示不排除任何记录）
func (r *categoryRepository) ExistsByName(ctx context.Context, name string, excludeID uint64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.Category{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// UpdateArticleCount 直接设置文章数量
func (r *categoryRepository) UpdateArticleCount(ctx context.Context, id uint64, count int) error {
	return r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", id).
		Update("article_count", count).Error
}

// IncrArticleCount 原子性增减文章计数（delta 可为负数）
func (r *categoryRepository) IncrArticleCount(ctx context.Context, id uint64, delta int) error {
	expr := "article_count + ?"
	if delta < 0 {
		// 防止计数变负
		return r.db.WithContext(ctx).
			Model(&model.Category{}).
			Where("id = ? AND article_count >= ?", id, -delta).
			Update("article_count", gorm.Expr("article_count + ?", delta)).Error
	}
	return r.db.WithContext(ctx).
		Model(&model.Category{}).
		Where("id = ?", id).
		Update("article_count", gorm.Expr(expr, delta)).Error
}

// ─── 模型映射 ────────────────────────────────────────────────────────────────

func categoryDomainToModel(c *domain.Category) *model.Category {
	return &model.Category{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		Cover:        c.Cover,
		SortOrder:    c.SortOrder,
		ArticleCount: c.ArticleCount,
	}
}

func categoryModelToDomain(m *model.Category) *domain.Category {
	return &domain.Category{
		ID:           m.ID,
		Name:         m.Name,
		Slug:         m.Slug,
		Description:  m.Description,
		Cover:        m.Cover,
		SortOrder:    m.SortOrder,
		ArticleCount: m.ArticleCount,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
