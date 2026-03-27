// Package model GORM 数据模型，与领域实体解耦，专注于持久化映射。
package model

import (
	"time"

	"gorm.io/gorm"
)

// User GORM 用户模型
type User struct {
	ID        uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	UUID      string         `gorm:"column:uuid;type:char(36);uniqueIndex;not null"`
	Username  string         `gorm:"column:username;type:varchar(32);uniqueIndex;not null"`
	Email     string         `gorm:"column:email;type:varchar(128);uniqueIndex;not null"`
	Password  string         `gorm:"column:password;type:varchar(255);not null"`
	Nickname  string         `gorm:"column:nickname;type:varchar(64);not null;default:''"`
	Avatar    string         `gorm:"column:avatar;type:varchar(512);not null;default:''"`
	Bio       string         `gorm:"column:bio;type:varchar(512);not null;default:''"`
	Role      int8           `gorm:"column:role;type:tinyint;not null;default:1"`
	Status    int8           `gorm:"column:status;type:tinyint;not null;default:1"`
	LastLogin *time.Time     `gorm:"column:last_login"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (User) TableName() string { return "users" }

// UserToken GORM 刷新令牌模型
type UserToken struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       uint64    `gorm:"column:user_id;not null;index"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar(512);not null;index"`
	UserAgent    string    `gorm:"column:user_agent;type:varchar(512);not null;default:''"`
	ClientIP     string    `gorm:"column:client_ip;type:varchar(64);not null;default:''"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null;index"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UserToken) TableName() string { return "user_tokens" }
