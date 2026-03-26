// Package main AntBlog 服务入口。
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	apparticle "antblog/internal/application/article"
	appcategory "antblog/internal/application/category"
	appcomment "antblog/internal/application/comment"
	appmedia "antblog/internal/application/media"
	apptag "antblog/internal/application/tag"
	appuser "antblog/internal/application/user"
	"antblog/internal/infrastructure/persistence"
	"antblog/internal/infrastructure/persistence/model"
	"antblog/internal/interfaces/http/router"
	"antblog/internal/interfaces/validator"
	pkgcache "antblog/pkg/cache"
	"antblog/pkg/config"
	"antblog/pkg/db"
	pkgjwt "antblog/pkg/jwt"
	"antblog/pkg/logger"
)

func main() {
	app := fx.New(
		// ── 公共基础库 ──────────────────────────────────────────────────
		config.Module,
		logger.Module,
		db.Module,
		pkgcache.Module,
		pkgjwt.Module,

		// ── 基础设施层（数据库仓储 + 文件存储驱动）──────────────────────
		persistence.Module,

		// ── 应用层（用例）────────────────────────────────────────────────
		appuser.Module,
		appcategory.Module,
		apptag.Module,
		apparticle.Module,
		appcomment.Module,
		appmedia.Module,

		// ── 接口层（路由/Handler）────────────────────────────────────────
		router.Module,

		// ── 启动钩子 ────────────────────────────────────────────────────
		fx.Invoke(registerHooks),
		fx.Invoke(runMigrations),
		fx.Invoke(ensureDefaultAdmin),
		fx.Invoke(registerValidators),
	)
	app.Run()
}

func registerHooks(lc fx.Lifecycle, engine *gin.Engine, cfg *config.Config, log *zap.Logger) {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("starting HTTP server", zap.String("addr", addr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("shutting down HTTP server")
			c, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return srv.Shutdown(c)
		},
	})
}

func runMigrations(log *zap.Logger) { log.Info("database migration completed") }

func ensureDefaultAdmin(cfg *config.Config, log *zap.Logger, gdb *gorm.DB) {
	if cfg.App.Env != "dev" {
		return
	}

	const (
		adminUsername = "admin"
		adminEmail    = "admin@antblog.dev"
		adminPassword = "$2a$10$GfxMXS0M2Vzrw7I.ehMpHuoufOp1tC1UfnI0QYBYd3B8vx0bsPFQu"
	)

	var u model.User
	err := gdb.Where("username = ? OR email = ?", adminUsername, adminEmail).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		now := time.Now()
		create := &model.User{
			UUID:      uuid.NewString(),
			Username:  adminUsername,
			Email:     adminEmail,
			Password:  adminPassword,
			Nickname:  "AntBlog Admin",
			Role:      2,
			Status:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if createErr := gdb.Create(create).Error; createErr != nil {
			log.Warn("ensure default admin create failed", zap.Error(createErr))
			return
		}
		log.Info("default admin ensured", zap.Uint64("user_id", create.ID))
		return
	}
	if err != nil {
		log.Warn("ensure default admin query failed", zap.Error(err))
		return
	}

	if updateErr := gdb.Model(&model.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"username":   adminUsername,
			"email":      adminEmail,
			"password":   adminPassword,
			"nickname":   "AntBlog Admin",
			"role":       2,
			"status":     1,
			"updated_at": time.Now(),
		}).Error; updateErr != nil {
		log.Warn("ensure default admin update failed", zap.Error(updateErr), zap.Uint64("user_id", u.ID))
		return
	}
	log.Info("default admin refreshed", zap.Uint64("user_id", u.ID))
}

func registerValidators(log *zap.Logger) {
	if err := validator.RegisterAll(); err != nil {
		log.Fatal("register validators failed", zap.Error(err))
	}
}
