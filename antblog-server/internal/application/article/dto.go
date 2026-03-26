// Package article 文章应用层。
package article

import (
	"time"

	domaintag "antblog/internal/application/tag"
)

// ─── 请求 DTO ────────────────────────────────────────────────────────────────

// CreateArticleReq 创建文章请求
type CreateArticleReq struct {
	Title        string   `json:"title"         validate:"required,max=256"`
	Slug         string   `json:"slug"          validate:"omitempty,max=300"`
	Summary      string   `json:"summary"       validate:"max=512"`
	Content      string   `json:"content"       validate:"required"`
	Cover        string   `json:"cover"         validate:"omitempty,max=512,url"`
	CategoryID   *uint64  `json:"category_id"`
	TagIDs       []uint64 `json:"tag_ids"`
	Status       int8     `json:"status"        validate:"omitempty,oneof=1 2 3"`
	IsTop        bool     `json:"is_top"`
	IsFeatured   bool     `json:"is_featured"`
	AllowComment bool     `json:"allow_comment"`
}

// UpdateArticleReq 更新文章请求
type UpdateArticleReq struct {
	Title        string   `json:"title"         validate:"required,max=256"`
	Slug         string   `json:"slug"          validate:"omitempty,max=300"`
	Summary      string   `json:"summary"       validate:"max=512"`
	Content      string   `json:"content"       validate:"required"`
	Cover        string   `json:"cover"         validate:"omitempty,max=512,url"`
	CategoryID   *uint64  `json:"category_id"`
	TagIDs       []uint64 `json:"tag_ids"`
	IsTop        bool     `json:"is_top"`
	IsFeatured   bool     `json:"is_featured"`
	AllowComment bool     `json:"allow_comment"`
}

// UpdateStatusReq 更新文章状态请求
type UpdateStatusReq struct {
	Status int8 `json:"status" validate:"required,oneof=1 2 3"`
}

// ListArticleReq 文章列表查询请求（前台）
type ListArticleReq struct {
	Page       int    `form:"page"        validate:"min=1"`
	PageSize   int    `form:"page_size"   validate:"min=1,max=50"`
	CategoryID uint64 `form:"category_id"`
	TagID      uint64 `form:"tag_id"`
	Keyword    string `form:"keyword"     validate:"max=100"`
}

// AdminListArticleReq 后台文章列表查询请求
type AdminListArticleReq struct {
	Page       int    `form:"page"        validate:"min=1"`
	PageSize   int    `form:"page_size"   validate:"min=1,max=50"`
	Status     int8   `form:"status"      validate:"omitempty,oneof=1 2 3"`
	CategoryID uint64 `form:"category_id"`
	TagID      uint64 `form:"tag_id"`
	Keyword    string `form:"keyword"     validate:"max=100"`
}

// ArchiveReq 归档时间线查询请求
type ArchiveDetailReq struct {
	Year     int `form:"year"      validate:"required,min=2000,max=2100"`
	Month    int `form:"month"     validate:"required,min=1,max=12"`
	Page     int `form:"page"      validate:"min=1"`
	PageSize int `form:"page_size" validate:"min=1,max=50"`
}

// ─── 响应 DTO ────────────────────────────────────────────────────────────────

// ArticleResp 文章完整响应（详情页用）
type ArticleResp struct {
	ID            uint64                   `json:"id"`
	UUID          string                   `json:"uuid"`
	AuthorID      uint64                   `json:"author_id"`
	CategoryID    *uint64                  `json:"category_id"`
	Title         string                   `json:"title"`
	Slug          string                   `json:"slug"`
	Summary       string                   `json:"summary"`
	Content       string                   `json:"content"`
	ContentHTML   string                   `json:"content_html"`
	Cover         string                   `json:"cover"`
	Status        int8                     `json:"status"`
	StatusText    string                   `json:"status_text"`
	IsTop         bool                     `json:"is_top"`
	IsFeatured    bool                     `json:"is_featured"`
	AllowComment  bool                     `json:"allow_comment"`
	ViewCount     int                      `json:"view_count"`
	LikeCount     int                      `json:"like_count"`
	CommentCount  int                      `json:"comment_count"`
	BookmarkCount int                      `json:"bookmark_count"`
	WordCount     int                      `json:"word_count"`
	PublishedAt   *time.Time               `json:"published_at"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
	Tags          []*domaintag.TagSimpleResp `json:"tags"`
	// 当前用户互动状态（需鉴权，未登录时为 nil）
	Liked      *bool `json:"liked,omitempty"`
	Bookmarked *bool `json:"bookmarked,omitempty"`
}

// ArticleListItemResp 文章列表条目响应（不含正文）
type ArticleListItemResp struct {
	ID            uint64                   `json:"id"`
	UUID          string                   `json:"uuid"`
	AuthorID      uint64                   `json:"author_id"`
	CategoryID    *uint64                  `json:"category_id"`
	Title         string                   `json:"title"`
	Slug          string                   `json:"slug"`
	Summary       string                   `json:"summary"`
	Cover         string                   `json:"cover"`
	Status        int8                     `json:"status"`
	IsTop         bool                     `json:"is_top"`
	IsFeatured    bool                     `json:"is_featured"`
	ViewCount     int                      `json:"view_count"`
	LikeCount     int                      `json:"like_count"`
	CommentCount  int                      `json:"comment_count"`
	BookmarkCount int                      `json:"bookmark_count"`
	WordCount     int                      `json:"word_count"`
	PublishedAt   *time.Time               `json:"published_at"`
	CreatedAt     time.Time                `json:"created_at"`
	Tags          []*domaintag.TagSimpleResp `json:"tags"`
	// 互动状态：仅登录用户访问时注入，未登录时为 null
	Liked      *bool `json:"liked,omitempty"`
	Bookmarked *bool `json:"bookmarked,omitempty"`
}

// ArchiveItemResp 归档时间线条目响应
type ArchiveItemResp struct {
	Year         int `json:"year"`
	Month        int `json:"month"`
	ArticleCount int `json:"article_count"`
}
