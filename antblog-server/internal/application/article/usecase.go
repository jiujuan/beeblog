package article

import "context"

// IArticleUseCase 文章用例接口
type IArticleUseCase interface {
	// ── 前台接口 ─────────────────────────────────────────────────────────────

	// ListArticles 前台文章列表（仅已发布，支持分类/标签/关键字过滤）
	// userID 可选，登录用户传入后列表项会注入 liked/bookmarked 状态
	ListArticles(ctx context.Context, req *ListArticleReq, userID *uint64) ([]*ArticleListItemResp, int64, error)

	// GetArticleBySlug 前台按 Slug 获取文章详情（阅读数 +1，注入互动状态）
	GetArticleBySlug(ctx context.Context, slug string, userID *uint64) (*ArticleResp, error)

	// GetFeaturedArticles 获取精选文章（首页推荐）
	// userID 可选，登录用户传入后会注入 liked/bookmarked 状态
	GetFeaturedArticles(ctx context.Context, limit int, userID *uint64) ([]*ArticleListItemResp, error)

	// GetArchive 获取归档时间线
	GetArchive(ctx context.Context) ([]*ArchiveItemResp, error)

	// GetArchiveDetail 获取某年月的文章列表
	// userID 可选
	GetArchiveDetail(ctx context.Context, req *ArchiveDetailReq, userID *uint64) ([]*ArticleListItemResp, int64, error)

	// ── 互动接口（需登录） ───────────────────────────────────────────────────

	// LikeArticle 点赞文章（幂等）
	LikeArticle(ctx context.Context, articleID, userID uint64) error

	// UnlikeArticle 取消点赞
	UnlikeArticle(ctx context.Context, articleID, userID uint64) error

	// BookmarkArticle 收藏文章（幂等）
	BookmarkArticle(ctx context.Context, articleID, userID uint64) error

	// UnbookmarkArticle 取消收藏
	UnbookmarkArticle(ctx context.Context, articleID, userID uint64) error

	// GetUserBookmarks 获取当前用户的收藏列表
	GetUserBookmarks(ctx context.Context, userID uint64, page, pageSize int) ([]*ArticleListItemResp, int64, error)

	// ── 后台管理接口 ──────────────────────────────────────────────────────────

	// AdminListArticles 后台文章列表（全状态，支持更多过滤维度）
	AdminListArticles(ctx context.Context, req *AdminListArticleReq) ([]*ArticleListItemResp, int64, error)

	// AdminGetArticle 后台按 ID 获取文章（含草稿）
	AdminGetArticle(ctx context.Context, id uint64) (*ArticleResp, error)

	// CreateArticle 创建文章
	CreateArticle(ctx context.Context, authorID uint64, req *CreateArticleReq) (*ArticleResp, error)

	// UpdateArticle 更新文章内容
	UpdateArticle(ctx context.Context, id uint64, req *UpdateArticleReq) (*ArticleResp, error)

	// UpdateArticleStatus 变更文章状态（发布/归档/撤稿）
	UpdateArticleStatus(ctx context.Context, id uint64, req *UpdateStatusReq) (*ArticleResp, error)

	// DeleteArticle 软删除文章
	DeleteArticle(ctx context.Context, id uint64) error
}
