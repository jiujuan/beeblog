package db

import (
	"fmt"

	"gorm.io/gorm"
)

// Migrator 数据库迁移管理器
type Migrator struct {
	db     *gorm.DB
	models []any
}

// NewMigrator 创建迁移管理器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// Register 注册需要迁移的 GORM 模型
func (m *Migrator) Register(models ...any) *Migrator {
	m.models = append(m.models, models...)
	return m
}

// AutoMigrate 执行自动迁移（创建/更新表结构，不删除列）
func (m *Migrator) AutoMigrate() error {
	if len(m.models) == 0 {
		return nil
	}
	if err := m.db.AutoMigrate(m.models...); err != nil {
		return fmt.Errorf("db: auto migrate: %w", err)
	}
	return nil
}

// CreateTable 仅创建表（若已存在则跳过）
func (m *Migrator) CreateTable(models ...any) error {
	for _, model := range models {
		if m.db.Migrator().HasTable(model) {
			continue
		}
		if err := m.db.Migrator().CreateTable(model); err != nil {
			return fmt.Errorf("db: create table for %T: %w", model, err)
		}
	}
	return nil
}

// DropTable 删除表（危险操作，仅用于测试）
func (m *Migrator) DropTable(models ...any) error {
	return m.db.Migrator().DropTable(models...)
}

// HasTable 检查表是否存在
func (m *Migrator) HasTable(model any) bool {
	return m.db.Migrator().HasTable(model)
}

// HasColumn 检查列是否存在
func (m *Migrator) HasColumn(model any, column string) bool {
	return m.db.Migrator().HasColumn(model, column)
}
