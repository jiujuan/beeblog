package comment

import "context"

// ICommentRepository 评论仓储接口（由 infrastructure 层实现）
type ICommentRepository interface {
	// ── 单条操作 ─────────────────────────────────────────────────────────────

	// Create 创建评论
	Create(ctx context.Context, c *Comment) (*Comment, error)

	// FindByID 按 ID 查询
	FindByID(ctx context.Context, id uint64) (*Comment, error)

	// Update 更新评论（审核状态等）
	Update(ctx context.Context, c *Comment) error

	// Delete 软删除评论
	Delete(ctx context.Context, id uint64) error

	// ── 前台列表 ─────────────────────────────────────────────────────────────

	// FindTopLevelByArticle 查询某文章的顶级已通过评论（分页，按时间升序）
	FindTopLevelByArticle(ctx context.Context, articleID uint64, page, pageSize int) ([]*Comment, int64, error)

	// FindChildrenByRoot 查询某根评论下的所有子评论（已通过，按时间升序）
	FindChildrenByRoot(ctx context.Context, rootID uint64) ([]*Comment, error)

	// ── 后台管理 ─────────────────────────────────────────────────────────────

	// AdminFind 后台评论列表（支持多条件过滤，分页）
	AdminFind(ctx context.Context, filter *AdminFilter) ([]*Comment, int64, error)

	// ── 计数 ─────────────────────────────────────────────────────────────────

	// CountByArticle 统计某文章已通过评论总数（用于文章 comment_count 同步）
	CountByArticle(ctx context.Context, articleID uint64) (int64, error)

	// IncrLikeCount 原子性增减点赞数
	IncrLikeCount(ctx context.Context, id uint64, delta int) error
}

// AdminFilter 后台评论过滤条件
type AdminFilter struct {
	ArticleID *uint64
	UserID    *uint64
	Status    *Status
	Keyword   string // 内容关键字模糊搜索
	Page      int
	PageSize  int
}
