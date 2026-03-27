package media

import "context"

// IMediaRepository 媒体仓储接口（由 infrastructure 层实现）
type IMediaRepository interface {
	// ── 单条操作 ─────────────────────────────────────────────────────────────

	// Create 创建媒体记录
	Create(ctx context.Context, m *Media) (*Media, error)

	// FindByID 按 ID 查询
	FindByID(ctx context.Context, id uint64) (*Media, error)

	// FindByHash 按文件哈希查询（去重用）
	FindByHash(ctx context.Context, hash string) (*Media, error)

	// Update 更新媒体记录（绑定/解绑文章）
	Update(ctx context.Context, m *Media) error

	// Delete 软删除媒体记录
	Delete(ctx context.Context, id uint64) error

	// ── 列表查询 ─────────────────────────────────────────────────────────────

	// FindByUploaderID 查询某用户上传的所有媒体（分页，按时间倒序）
	FindByUploaderID(ctx context.Context, uploaderID uint64, page, pageSize int) ([]*Media, int64, error)

	// FindByArticleID 查询某文章关联的所有媒体
	FindByArticleID(ctx context.Context, articleID uint64) ([]*Media, error)

	// AdminFind 后台列表（支持多条件过滤，分页）
	AdminFind(ctx context.Context, filter *AdminFilter) ([]*Media, int64, error)

	// ── 批量操作 ─────────────────────────────────────────────────────────────

	// DeleteByUploaderID 删除某用户的所有媒体（账户注销时使用）
	DeleteByUploaderID(ctx context.Context, uploaderID uint64) error
}

// AdminFilter 后台媒体列表过滤条件
type AdminFilter struct {
	UploaderID *uint64
	ArticleID  *uint64
	MimeType   string // 如 "image/jpeg"，空字符串表示不过滤
	Page       int
	PageSize   int
}
