package article

import (
	"context"
	"time"
)

// IArticleRepository 文章仓储接口（由 infrastructure 层实现）
type IArticleRepository interface {
	// ── 单条操作 ─────────────────────────────────────────────────────────────

	// Create 创建文章（含标签关联写入）
	Create(ctx context.Context, a *Article) (*Article, error)

	// FindByID 按 ID 查询（不区分状态，供后台使用）
	FindByID(ctx context.Context, id uint64) (*Article, error)

	// FindBySlug 按 Slug 查询（前台使用，仅返回已发布文章）
	FindBySlug(ctx context.Context, slug string) (*Article, error)

	// FindByUUID 按 UUID 查询
	FindByUUID(ctx context.Context, uuid string) (*Article, error)

	// Update 更新文章主体字段
	Update(ctx context.Context, a *Article) error

	// UpdateStatus 单独更新文章状态（含 published_at 处理）
	UpdateStatus(ctx context.Context, id uint64, status Status, publishedAt *time.Time) error

	// Delete 软删除文章
	Delete(ctx context.Context, id uint64) error

	// ── 列表查询 ─────────────────────────────────────────────────────────────

	// FindList 通用分页列表（前台/后台共用，通过 filter 区分）
	FindList(ctx context.Context, filter *ListFilter) ([]*Article, int64, error)

	// FindFeatured 获取精选文章（前台首页）
	FindFeatured(ctx context.Context, limit int) ([]*Article, error)

	// FindByTagID 按标签分页（前台标签页）
	FindByTagID(ctx context.Context, tagID uint64, page, pageSize int) ([]*Article, int64, error)

	// FindByCategoryID 按分类分页（前台分类页）
	FindByCategoryID(ctx context.Context, categoryID uint64, page, pageSize int) ([]*Article, int64, error)

	// ── 归档 ─────────────────────────────────────────────────────────────────

	// GetArchive 按年月聚合归档数据
	GetArchive(ctx context.Context) ([]*ArchiveItem, error)

	// FindByYearMonth 按年月分页查文章（归档详情页）
	FindByYearMonth(ctx context.Context, year, month, page, pageSize int) ([]*Article, int64, error)

	// ── 标签关联 ─────────────────────────────────────────────────────────────

	// SyncTags 同步文章-标签关联（先删后写，事务内）
	SyncTags(ctx context.Context, articleID uint64, tagIDs []uint64) error

	// FindTagIDsByArticleID 查询文章绑定的所有标签 ID
	FindTagIDsByArticleID(ctx context.Context, articleID uint64) ([]uint64, error)

	// ── 计数 ─────────────────────────────────────────────────────────────────

	// IncrViewCount 原子递增阅读数
	IncrViewCount(ctx context.Context, id uint64) error

	// IncrLikeCount 原子增减点赞数（delta=+1 或 -1）
	IncrLikeCount(ctx context.Context, id uint64, delta int) error

	// IncrCommentCount 原子增减评论数
	IncrCommentCount(ctx context.Context, id uint64, delta int) error

	// IncrBookmarkCount 原子增减收藏数
	IncrBookmarkCount(ctx context.Context, id uint64, delta int) error

	// ── 互动状态（点赞/收藏） ─────────────────────────────────────────────────

	// AddLike 点赞（写 article_likes，幂等）
	AddLike(ctx context.Context, articleID, userID uint64) error

	// RemoveLike 取消点赞
	RemoveLike(ctx context.Context, articleID, userID uint64) error

	// HasLiked 查询用户是否已点赞
	HasLiked(ctx context.Context, articleID, userID uint64) (bool, error)

	// AddBookmark 收藏（写 article_bookmarks，幂等）
	AddBookmark(ctx context.Context, articleID, userID uint64) error

	// RemoveBookmark 取消收藏
	RemoveBookmark(ctx context.Context, articleID, userID uint64) error

	// HasBookmarked 查询用户是否已收藏
	HasBookmarked(ctx context.Context, articleID, userID uint64) (bool, error)

	// BatchHasLiked 批量查询用户对一批文章的点赞状态，返回已点赞的 articleID 集合
	// 用于列表渲染，避免 N+1 查询
	BatchHasLiked(ctx context.Context, articleIDs []uint64, userID uint64) (map[uint64]bool, error)

	// BatchHasBookmarked 批量查询用户对一批文章的收藏状态，返回已收藏的 articleID 集合
	BatchHasBookmarked(ctx context.Context, articleIDs []uint64, userID uint64) (map[uint64]bool, error)

	// GetUserBookmarks 获取用户收藏列表（分页）
	GetUserBookmarks(ctx context.Context, userID uint64, page, pageSize int) ([]*Article, int64, error)

	// ── 唯一性校验 ───────────────────────────────────────────────────────────

	// ExistsBySlug 检查 Slug 是否存在（excludeID=0 表示不排除）
	ExistsBySlug(ctx context.Context, slug string, excludeID uint64) (bool, error)
}
