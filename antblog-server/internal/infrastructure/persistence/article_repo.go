package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "antblog/internal/domain/article"
	"antblog/internal/infrastructure/persistence/model"
	apperrors "antblog/pkg/errors"
)

// articleRepository IArticleRepository 的 GORM 实现
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储
func NewArticleRepository(db *gorm.DB) domain.IArticleRepository {
	return &articleRepository{db: db}
}

// ─── 单条操作 ────────────────────────────────────────────────────────────────

// Create 创建文章（含标签关联，事务内完成）
func (r *articleRepository) Create(ctx context.Context, a *domain.Article) (*domain.Article, error) {
	m := articleDomainToModel(a)
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		if len(a.TagIDs) > 0 {
			return syncTagsTx(tx, m.ID, a.TagIDs)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	result := articleModelToDomain(m)
	result.TagIDs = a.TagIDs
	return result, nil
}

// FindByID 按 ID 查询（含软删除过滤，含标签 ID 列表）
func (r *articleRepository) FindByID(ctx context.Context, id uint64) (*domain.Article, error) {
	var m model.Article
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrArticleNotFound()
		}
		return nil, err
	}
	d := articleModelToDomain(&m)
	d.TagIDs, _ = r.FindTagIDsByArticleID(ctx, id)
	return d, nil
}

// FindBySlug 按 Slug 查询（仅已发布）
func (r *articleRepository) FindBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	var m model.Article
	err := r.db.WithContext(ctx).
		Where("slug = ? AND status = ?", slug, domain.StatusPublished).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrArticleNotFound()
		}
		return nil, err
	}
	d := articleModelToDomain(&m)
	d.TagIDs, _ = r.FindTagIDsByArticleID(ctx, m.ID)
	return d, nil
}

// FindByUUID 按 UUID 查询
func (r *articleRepository) FindByUUID(ctx context.Context, uuid string) (*domain.Article, error) {
	var m model.Article
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrArticleNotFound()
		}
		return nil, err
	}
	d := articleModelToDomain(&m)
	d.TagIDs, _ = r.FindTagIDsByArticleID(ctx, m.ID)
	return d, nil
}

// Update 更新文章主体字段（含标签同步，事务内）
func (r *articleRepository) Update(ctx context.Context, a *domain.Article) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"category_id":   a.CategoryID,
			"title":         a.Title,
			"slug":          a.Slug,
			"summary":       a.Summary,
			"content":       a.Content,
			"content_html":  a.ContentHTML,
			"cover":         a.Cover,
			"is_top":        boolToInt8(a.IsTop),
			"is_featured":   boolToInt8(a.IsFeatured),
			"allow_comment": boolToInt8(a.AllowComment),
			"word_count":    a.WordCount,
		}
		if err := tx.Model(&model.Article{}).Where("id = ?", a.ID).Updates(updates).Error; err != nil {
			return err
		}
		return syncTagsTx(tx, a.ID, a.TagIDs)
	})
}

// UpdateStatus 单独更新状态
func (r *articleRepository) UpdateStatus(ctx context.Context, id uint64, status domain.Status, publishedAt *time.Time) error {
	updates := map[string]any{"status": int8(status)}
	if publishedAt != nil {
		updates["published_at"] = publishedAt
	}
	return r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 软删除
func (r *articleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Article{}, "id = ?", id).Error
}

// ─── 列表查询 ────────────────────────────────────────────────────────────────

// FindList 通用分页列表
func (r *articleRepository) FindList(ctx context.Context, filter *domain.ListFilter) ([]*domain.Article, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Article{})

	if filter.Status != nil {
		q = q.Where("status = ?", int8(*filter.Status))
	}
	if filter.CategoryID != nil {
		q = q.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.AuthorID != nil {
		q = q.Where("author_id = ?", *filter.AuthorID)
	}
	if filter.IsTop != nil {
		q = q.Where("is_top = ?", boolToInt8(*filter.IsTop))
	}
	if filter.IsFeatured != nil {
		q = q.Where("is_featured = ?", boolToInt8(*filter.IsFeatured))
	}
	if filter.Keyword != "" {
		q = q.Where("title LIKE ?", "%"+filter.Keyword+"%")
	}
	// 按标签过滤（子查询）
	if filter.TagID != nil {
		q = q.Where("id IN (?)",
			r.db.Table("article_tags").Select("article_id").Where("tag_id = ?", *filter.TagID),
		)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}

	var list []model.Article
	offset := (filter.Page - 1) * filter.PageSize
	if err := q.
		Order("is_top DESC, published_at DESC, created_at DESC").
		Offset(offset).Limit(filter.PageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return r.modelsToDomainsWithTags(ctx, list), total, nil
}

// FindFeatured 获取精选文章
func (r *articleRepository) FindFeatured(ctx context.Context, limit int) ([]*domain.Article, error) {
	var list []model.Article
	err := r.db.WithContext(ctx).
		Where("status = ? AND is_featured = ?", domain.StatusPublished, 1).
		Order("published_at DESC").
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToDomainsWithTags(ctx, list), nil
}

// FindByTagID 按标签分页
func (r *articleRepository) FindByTagID(ctx context.Context, tagID uint64, page, pageSize int) ([]*domain.Article, int64, error) {
	status := domain.StatusPublished
	filter := &domain.ListFilter{
		Status:   &status,
		TagID:    &tagID,
		Page:     page,
		PageSize: pageSize,
	}
	return r.FindList(ctx, filter)
}

// FindByCategoryID 按分类分页
func (r *articleRepository) FindByCategoryID(ctx context.Context, categoryID uint64, page, pageSize int) ([]*domain.Article, int64, error) {
	status := domain.StatusPublished
	filter := &domain.ListFilter{
		Status:     &status,
		CategoryID: &categoryID,
		Page:       page,
		PageSize:   pageSize,
	}
	return r.FindList(ctx, filter)
}

// ─── 归档 ────────────────────────────────────────────────────────────────────

// GetArchive 按年月聚合归档（直接走 DB，可用视图）
func (r *articleRepository) GetArchive(ctx context.Context) ([]*domain.ArchiveItem, error) {
	type archiveRow struct {
		Year         int `gorm:"column:year"`
		Month        int `gorm:"column:month"`
		ArticleCount int `gorm:"column:article_count"`
	}
	var rows []archiveRow
	err := r.db.WithContext(ctx).
		Table("articles").
		Select("YEAR(published_at) AS year, MONTH(published_at) AS month, COUNT(*) AS article_count").
		Where("status = ? AND deleted_at IS NULL AND published_at IS NOT NULL", int8(domain.StatusPublished)).
		Group("YEAR(published_at), MONTH(published_at)").
		Order("year DESC, month DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]*domain.ArchiveItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, &domain.ArchiveItem{
			Year:         row.Year,
			Month:        row.Month,
			ArticleCount: row.ArticleCount,
		})
	}
	return result, nil
}

// FindByYearMonth 按年月分页
func (r *articleRepository) FindByYearMonth(ctx context.Context, year, month, page, pageSize int) ([]*domain.Article, int64, error) {
	var list []model.Article
	q := r.db.WithContext(ctx).
		Where("status = ? AND YEAR(published_at) = ? AND MONTH(published_at) = ?",
			int8(domain.StatusPublished), year, month)

	var total int64
	if err := q.Model(&model.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := q.Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return r.modelsToDomainsWithTags(ctx, list), total, nil
}

// ─── 标签关联 ────────────────────────────────────────────────────────────────

// SyncTags 同步文章标签关联（事务外调用入口）
func (r *articleRepository) SyncTags(ctx context.Context, articleID uint64, tagIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return syncTagsTx(tx, articleID, tagIDs)
	})
}

// FindTagIDsByArticleID 查询文章的标签 ID 列表
func (r *articleRepository) FindTagIDsByArticleID(ctx context.Context, articleID uint64) ([]uint64, error) {
	var tagIDs []uint64
	err := r.db.WithContext(ctx).
		Table("article_tags").
		Select("tag_id").
		Where("article_id = ?", articleID).
		Pluck("tag_id", &tagIDs).Error
	return tagIDs, err
}

// ─── 计数 ────────────────────────────────────────────────────────────────────

func (r *articleRepository) IncrViewCount(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *articleRepository) IncrLikeCount(ctx context.Context, id uint64, delta int) error {
	return incrCountSafe(r.db.WithContext(ctx), id, "like_count", delta)
}

func (r *articleRepository) IncrCommentCount(ctx context.Context, id uint64, delta int) error {
	return incrCountSafe(r.db.WithContext(ctx), id, "comment_count", delta)
}

func (r *articleRepository) IncrBookmarkCount(ctx context.Context, id uint64, delta int) error {
	return incrCountSafe(r.db.WithContext(ctx), id, "bookmark_count", delta)
}

// ─── 互动：点赞 ──────────────────────────────────────────────────────────────

func (r *articleRepository) AddLike(ctx context.Context, articleID, userID uint64) error {
	like := &model.ArticleLike{ArticleID: articleID, UserID: userID}
	return r.db.WithContext(ctx).
		Where(model.ArticleLike{ArticleID: articleID, UserID: userID}).
		FirstOrCreate(like).Error
}

func (r *articleRepository) RemoveLike(ctx context.Context, articleID, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Delete(&model.ArticleLike{}).Error
}

func (r *articleRepository) HasLiked(ctx context.Context, articleID, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ArticleLike{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count).Error
	return count > 0, err
}

// ─── 互动：收藏 ──────────────────────────────────────────────────────────────

func (r *articleRepository) AddBookmark(ctx context.Context, articleID, userID uint64) error {
	bm := &model.ArticleBookmark{ArticleID: articleID, UserID: userID}
	return r.db.WithContext(ctx).
		Where(model.ArticleBookmark{ArticleID: articleID, UserID: userID}).
		FirstOrCreate(bm).Error
}

func (r *articleRepository) RemoveBookmark(ctx context.Context, articleID, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Delete(&model.ArticleBookmark{}).Error
}

func (r *articleRepository) HasBookmarked(ctx context.Context, articleID, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ArticleBookmark{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count).Error
	return count > 0, err
}

// BatchHasLiked 批量查询用户是否点赞了一批文章，返回已点赞的 articleID → true 映射
func (r *articleRepository) BatchHasLiked(ctx context.Context, articleIDs []uint64, userID uint64) (map[uint64]bool, error) {
	result := make(map[uint64]bool, len(articleIDs))
	if len(articleIDs) == 0 || userID == 0 {
		return result, nil
	}
	var liked []uint64
	err := r.db.WithContext(ctx).
		Model(&model.ArticleLike{}).
		Select("article_id").
		Where("article_id IN ? AND user_id = ?", articleIDs, userID).
		Pluck("article_id", &liked).Error
	if err != nil {
		return result, err
	}
	for _, id := range liked {
		result[id] = true
	}
	return result, nil
}

// BatchHasBookmarked 批量查询用户是否收藏了一批文章，返回已收藏的 articleID → true 映射
func (r *articleRepository) BatchHasBookmarked(ctx context.Context, articleIDs []uint64, userID uint64) (map[uint64]bool, error) {
	result := make(map[uint64]bool, len(articleIDs))
	if len(articleIDs) == 0 || userID == 0 {
		return result, nil
	}
	var bookmarked []uint64
	err := r.db.WithContext(ctx).
		Model(&model.ArticleBookmark{}).
		Select("article_id").
		Where("article_id IN ? AND user_id = ?", articleIDs, userID).
		Pluck("article_id", &bookmarked).Error
	if err != nil {
		return result, err
	}
	for _, id := range bookmarked {
		result[id] = true
	}
	return result, nil
}

// GetUserBookmarks 获取用户收藏的文章列表
func (r *articleRepository) GetUserBookmarks(ctx context.Context, userID uint64, page, pageSize int) ([]*domain.Article, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.ArticleBookmark{}).
		Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var bookmarks []model.ArticleBookmark
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&bookmarks).Error; err != nil {
		return nil, 0, err
	}

	articleIDs := make([]uint64, 0, len(bookmarks))
	for _, bm := range bookmarks {
		articleIDs = append(articleIDs, bm.ArticleID)
	}
	if len(articleIDs) == 0 {
		return nil, total, nil
	}

	var list []model.Article
	if err := r.db.WithContext(ctx).Where("id IN ?", articleIDs).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return r.modelsToDomainsWithTags(ctx, list), total, nil
}

// ─── 唯一性校验 ──────────────────────────────────────────────────────────────

func (r *articleRepository) ExistsBySlug(ctx context.Context, slug string, excludeID uint64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.Article{}).Where("slug = ?", slug)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// ─── 私有辅助 ────────────────────────────────────────────────────────────────

// syncTagsTx 在事务内同步文章-标签关联（先删后写）
func syncTagsTx(tx *gorm.DB, articleID uint64, tagIDs []uint64) error {
	if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleTag{}).Error; err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}
	ats := make([]model.ArticleTag, 0, len(tagIDs))
	for _, tid := range tagIDs {
		ats = append(ats, model.ArticleTag{ArticleID: articleID, TagID: tid})
	}
	return tx.CreateInBatches(ats, 50).Error
}

// incrCountSafe 原子增减计数字段（delta 为负时防止变负）
func incrCountSafe(db *gorm.DB, id uint64, col string, delta int) error {
	if delta < 0 {
		return db.Model(&model.Article{}).
			Where("id = ? AND "+col+" >= ?", id, -delta).
			UpdateColumn(col, gorm.Expr(col+" + ?", delta)).Error
	}
	return db.Model(&model.Article{}).
		Where("id = ?", id).
		UpdateColumn(col, gorm.Expr(col+" + ?", delta)).Error
}

// boolToInt8 将 bool 转换为数据库存储的 int8
func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// ─── 批量取标签（内部辅助）───────────────────────────────────────────────────

// modelsToDomainsWithTags 批量转换 model → domain，并一次性查询所有标签 ID
func (r *articleRepository) modelsToDomainsWithTags(ctx context.Context, list []model.Article) []*domain.Article {
	if len(list) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(list))
	for _, m := range list {
		ids = append(ids, m.ID)
	}

	// 一次查询所有文章的标签 ID，减少 N+1
	type row struct {
		ArticleID uint64 `gorm:"column:article_id"`
		TagID     uint64 `gorm:"column:tag_id"`
	}
	var rows []row
	_ = r.db.WithContext(ctx).
		Table("article_tags").
		Select("article_id, tag_id").
		Where("article_id IN ?", ids).
		Scan(&rows)

	tagMap := make(map[uint64][]uint64, len(list))
	for _, row := range rows {
		tagMap[row.ArticleID] = append(tagMap[row.ArticleID], row.TagID)
	}

	result := make([]*domain.Article, 0, len(list))
	for i := range list {
		d := articleModelToDomain(&list[i])
		d.TagIDs = tagMap[d.ID]
		result = append(result, d)
	}
	return result
}

// ─── 模型映射 ────────────────────────────────────────────────────────────────

func articleDomainToModel(a *domain.Article) *model.Article {
	return &model.Article{
		ID:            a.ID,
		UUID:          a.UUID,
		AuthorID:      a.AuthorID,
		CategoryID:    a.CategoryID,
		Title:         a.Title,
		Slug:          a.Slug,
		Summary:       a.Summary,
		Content:       a.Content,
		ContentHTML:   a.ContentHTML,
		Cover:         a.Cover,
		Status:        int8(a.Status),
		IsTop:         boolToInt8(a.IsTop),
		IsFeatured:    boolToInt8(a.IsFeatured),
		AllowComment:  boolToInt8(a.AllowComment),
		ViewCount:     a.ViewCount,
		LikeCount:     a.LikeCount,
		CommentCount:  a.CommentCount,
		BookmarkCount: a.BookmarkCount,
		WordCount:     a.WordCount,
		PublishedAt:   a.PublishedAt,
	}
}

func articleModelToDomain(m *model.Article) *domain.Article {
	return &domain.Article{
		ID:            m.ID,
		UUID:          m.UUID,
		AuthorID:      m.AuthorID,
		CategoryID:    m.CategoryID,
		Title:         m.Title,
		Slug:          m.Slug,
		Summary:       m.Summary,
		Content:       m.Content,
		ContentHTML:   m.ContentHTML,
		Cover:         m.Cover,
		Status:        domain.Status(m.Status),
		IsTop:         m.IsTop == 1,
		IsFeatured:    m.IsFeatured == 1,
		AllowComment:  m.AllowComment == 1,
		ViewCount:     m.ViewCount,
		LikeCount:     m.LikeCount,
		CommentCount:  m.CommentCount,
		BookmarkCount: m.BookmarkCount,
		WordCount:     m.WordCount,
		PublishedAt:   m.PublishedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
