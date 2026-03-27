package persistence

import (
	"go.uber.org/fx"

	domainarticle  "antblog/internal/domain/article"
	domaincategory "antblog/internal/domain/category"
	domaincomment  "antblog/internal/domain/comment"
	domainmedia    "antblog/internal/domain/media"
	domaintag      "antblog/internal/domain/tag"
	domainuser     "antblog/internal/domain/user"
	"antblog/internal/infrastructure/storage"
)

// Module fx 持久化模块，将所有 Repository 实现与存储驱动绑定为接口
var Module = fx.Options(
	fx.Provide(
		fx.Annotate(NewUserRepository,     fx.As(new(domainuser.IUserRepository))),
		fx.Annotate(NewCategoryRepository, fx.As(new(domaincategory.ICategoryRepository))),
		fx.Annotate(NewTagRepository,      fx.As(new(domaintag.ITagRepository))),
		fx.Annotate(NewArticleRepository,  fx.As(new(domainarticle.IArticleRepository))),
		fx.Annotate(NewCommentRepository,  fx.As(new(domaincomment.ICommentRepository))),
		fx.Annotate(NewMediaRepository,    fx.As(new(domainmedia.IMediaRepository))),
	),
	// 存储驱动：默认本地存储，可替换为 OSS
	fx.Provide(
		fx.Annotate(storage.NewLocalStorage, fx.As(new(domainmedia.IStorageDriver))),
	),
)
