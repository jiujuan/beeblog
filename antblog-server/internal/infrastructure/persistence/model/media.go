package model

import (
	"time"

	"gorm.io/gorm"
)

// Media GORM 媒体资源模型
type Media struct {
	ID           uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	UploaderID   uint64         `gorm:"column:uploader_id;not null;index:idx_media_uploader"`
	ArticleID    *uint64        `gorm:"column:article_id;index:idx_media_article"`
	OriginalName string         `gorm:"column:original_name;type:varchar(256);not null;default:''"`
	StoragePath  string         `gorm:"column:storage_path;type:varchar(512);not null"`
	URL          string         `gorm:"column:url;type:varchar(512);not null"`
	MimeType     string         `gorm:"column:mime_type;type:varchar(128);not null;default:'';index:idx_media_mime"`
	FileSize     int64          `gorm:"column:file_size;not null;default:0"`
	Width        int            `gorm:"column:width;not null;default:0"`
	Height       int            `gorm:"column:height;not null;default:0"`
	Hash         string         `gorm:"column:hash;type:char(64);not null;default:'';index:idx_media_hash"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index:idx_media_deleted"`
}

func (Media) TableName() string { return "media" }
