// Package persistence 基础设施层持久化实现，通过 GORM 操作数据库，
// 满足领域层定义的 IUserRepository 接口。
package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "antblog/internal/domain/user"
	"antblog/internal/infrastructure/persistence/model"
	apperrors "antblog/pkg/errors"
)

// userRepository IUserRepository 的 GORM 实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) domain.IUserRepository {
	return &userRepository{db: db}
}

// ─── 用户 CRUD ───────────────────────────────────────────────────────────────

func (r *userRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	m := domainToModel(u)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return modelToDomain(m), nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound()
		}
		return nil, err
	}
	return modelToDomain(&m), nil
}

func (r *userRepository) FindByUUID(ctx context.Context, uuid string) (*domain.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound()
		}
		return nil, err
	}
	return modelToDomain(&m), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound()
		}
		return nil, err
	}
	return modelToDomain(&m), nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound()
		}
		return nil, err
	}
	return modelToDomain(&m), nil
}

func (r *userRepository) Update(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"nickname":   u.Nickname,
			"avatar":     u.Avatar,
			"bio":        u.Bio,
			"password":   u.Password,
			"status":     u.Status,
			"role":       u.Role,
			"updated_at": time.Now(),
		}).Error
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("last_login", now).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ─── Token 管理 ──────────────────────────────────────────────────────────────

func (r *userRepository) SaveToken(ctx context.Context, token *domain.UserToken) error {
	m := &model.UserToken{
		UserID:       token.UserID,
		RefreshToken: token.RefreshToken,
		UserAgent:    token.UserAgent,
		ClientIP:     token.ClientIP,
		ExpiresAt:    token.ExpiresAt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *userRepository) FindToken(ctx context.Context, refreshToken string) (*domain.UserToken, error) {
	var m model.UserToken
	err := r.db.WithContext(ctx).
		Where("refresh_token = ?", refreshToken).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domain.UserToken{
		ID:           m.ID,
		UserID:       m.UserID,
		RefreshToken: m.RefreshToken,
		UserAgent:    m.UserAgent,
		ClientIP:     m.ClientIP,
		ExpiresAt:    m.ExpiresAt,
		CreatedAt:    m.CreatedAt,
	}, nil
}

func (r *userRepository) DeleteToken(ctx context.Context, refreshToken string) error {
	return r.db.WithContext(ctx).
		Where("refresh_token = ?", refreshToken).
		Delete(&model.UserToken{}).Error
}

func (r *userRepository) DeleteUserTokens(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.UserToken{}).Error
}

func (r *userRepository) CleanExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&model.UserToken{}).Error
}

// ─── 模型映射 ────────────────────────────────────────────────────────────────

// domainToModel 领域实体 → GORM 模型
func domainToModel(u *domain.User) *model.User {
	return &model.User{
		ID:        u.ID,
		UUID:      u.UUID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Bio:       u.Bio,
		Role:      int8(u.Role),
		Status:    int8(u.Status),
		LastLogin: u.LastLogin,
	}
}

// modelToDomain GORM 模型 → 领域实体
func modelToDomain(m *model.User) *domain.User {
	return &domain.User{
		ID:        m.ID,
		UUID:      m.UUID,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		Bio:       m.Bio,
		Role:      domain.Role(m.Role),
		Status:    domain.Status(m.Status),
		LastLogin: m.LastLogin,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
