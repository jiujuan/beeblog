package model

import (
	"time"

	"gorm.io/gorm"
)

// Article GORM 文章主表模型
type Article struct {
	ID            uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	UUID          string         `gorm:"column:uuid;type:char(36);not null;uniqueIndex"`
	AuthorID      uint64         `gorm:"column:author_id;not null;index"`
	CategoryID    *uint64        `gorm:"column:category_id;index"`
	Title         string         `gorm:"column:title;type:varchar(256);not null"`
	Slug          string         `gorm:"column:slug;type:varchar(300);not null;uniqueIndex"`
	Summary       string         `gorm:"column:summary;type:varchar(512);not null;default:''"`
	Content       string         `gorm:"column:content;type:longtext;not null"`
	ContentHTML   string         `gorm:"column:content_html;type:longtext;not null"`
	Cover         string         `gorm:"column:cover;type:varchar(512);not null;default:''"`
	Status        int8           `gorm:"column:status;not null;default:1;index"`
	IsTop         int8           `gorm:"column:is_top;not null;default:0"`
	IsFeatured    int8           `gorm:"column:is_featured;not null;default:0"`
	AllowComment  int8           `gorm:"column:allow_comment;not null;default:1"`
	ViewCount     int            `gorm:"column:view_count;not null;default:0"`
	LikeCount     int            `gorm:"column:like_count;not null;default:0"`
	CommentCount  int            `gorm:"column:comment_count;not null;default:0"`
	BookmarkCount int            `gorm:"column:bookmark_count;not null;default:0"`
	WordCount     int            `gorm:"column:word_count;not null;default:0"`
	PublishedAt   *time.Time     `gorm:"column:published_at;index"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Article) TableName() string { return "articles" }

// ArticleLike GORM 点赞模型
type ArticleLike struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	ArticleID uint64    `gorm:"column:article_id;not null;index:uq_likes,unique,priority:1"`
	UserID    uint64    `gorm:"column:user_id;not null;index:uq_likes,unique,priority:2;index:idx_likes_user"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ArticleLike) TableName() string { return "article_likes" }

// ArticleBookmark GORM 收藏模型
type ArticleBookmark struct {
	ID        uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	ArticleID uint64    `gorm:"column:article_id;not null;index:uq_bookmarks,unique,priority:1"`
	UserID    uint64    `gorm:"column:user_id;not null;index:uq_bookmarks,unique,priority:2;index:idx_bookmarks_user"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (ArticleBookmark) TableName() string { return "article_bookmarks" }
