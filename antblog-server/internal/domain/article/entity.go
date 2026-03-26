// Package article 文章领域层 —— 纯业务实体，零外部依赖。
package article

import "time"

// ─── 枚举 ────────────────────────────────────────────────────────────────────

// Status 文章发布状态
type Status int8

const (
	StatusDraft     Status = 1 // 草稿
	StatusPublished Status = 2 // 已发布
	StatusArchived  Status = 3 // 已归档
)

func (s Status) String() string {
	switch s {
	case StatusDraft:
		return "draft"
	case StatusPublished:
		return "published"
	case StatusArchived:
		return "archived"
	default:
		return "unknown"
	}
}

// IsValid 校验状态值是否合法
func (s Status) IsValid() bool {
	return s == StatusDraft || s == StatusPublished || s == StatusArchived
}

// ─── 聚合根 ──────────────────────────────────────────────────────────────────

// Article 文章聚合根
type Article struct {
	ID            uint64
	UUID          string
	AuthorID      uint64
	CategoryID    *uint64    // nil = 未分类
	Title         string
	Slug          string
	Summary       string     // 摘要；为空时可自动截取正文前 200 字
	Content       string     // Markdown 原文
	ContentHTML   string     // 服务端预渲染 HTML
	Cover         string     // 封面图 URL
	Status        Status
	IsTop         bool
	IsFeatured    bool
	AllowComment  bool
	ViewCount     int
	LikeCount     int
	CommentCount  int
	BookmarkCount int
	WordCount     int
	PublishedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// 关联标签 ID 列表（领域内聚合使用，存储在 article_tags 表）
	TagIDs []uint64
}

// ─── 业务方法 ────────────────────────────────────────────────────────────────

// Publish 发布文章（仅在草稿/归档状态下可操作）
// 首次发布时记录 published_at，重新发布不覆盖原始时间。
func (a *Article) Publish() {
	if a.Status != StatusPublished {
		a.Status = StatusPublished
		now := time.Now()
		if a.PublishedAt == nil {
			a.PublishedAt = &now
		}
		a.UpdatedAt = now
	}
}

// Archive 归档文章
func (a *Article) Archive() {
	a.Status = StatusArchived
	a.UpdatedAt = time.Now()
}

// Unpublish 撤回为草稿
func (a *Article) Unpublish() {
	a.Status = StatusDraft
	a.UpdatedAt = time.Now()
}

// IsPublished 是否已发布
func (a *Article) IsPublished() bool { return a.Status == StatusPublished }

// IsVisibleToPublic 前台是否可见（仅已发布文章可见）
func (a *Article) IsVisibleToPublic() bool { return a.Status == StatusPublished }

// SetTop 设置/取消置顶
func (a *Article) SetTop(v bool) {
	a.IsTop = v
	a.UpdatedAt = time.Now()
}

// SetFeatured 设置/取消精选
func (a *Article) SetFeatured(v bool) {
	a.IsFeatured = v
	a.UpdatedAt = time.Now()
}

// UpdateContent 更新正文及元数据（编辑保存时调用）
func (a *Article) UpdateContent(
	title, slug, summary, content, contentHTML, cover string,
	categoryID *uint64,
	tagIDs []uint64,
	allowComment, isTop, isFeatured bool,
) {
	a.Title = title
	a.Slug = slug
	a.Summary = summary
	a.Content = content
	a.ContentHTML = contentHTML
	a.Cover = cover
	a.CategoryID = categoryID
	a.TagIDs = tagIDs
	a.AllowComment = allowComment
	a.IsTop = isTop
	a.IsFeatured = isFeatured
	a.UpdatedAt = time.Now()
}

// IncrViewCount 阅读数 +1（供仓储层原子递增）
func (a *Article) IncrViewCount() { a.ViewCount++ }

// ─── 附属值对象 ──────────────────────────────────────────────────────────────

// ArchiveItem 时间线归档条目
type ArchiveItem struct {
	Year         int
	Month        int
	ArticleCount int
}

// ListFilter 文章列表查询过滤参数
type ListFilter struct {
	Status     *Status // nil = 不过滤
	CategoryID *uint64
	TagID      *uint64
	AuthorID   *uint64
	Keyword    string  // 标题模糊搜索
	IsTop      *bool
	IsFeatured *bool
	Page       int
	PageSize   int
}
