package article

import (
	"context"
	"time"

	"github.com/google/uuid"

	apperrors "antblog/pkg/errors"
	"antblog/pkg/utils"
)

// IDomainService 文章领域服务接口
type IDomainService interface {
	// ValidateCreate 校验创建参数
	ValidateCreate(ctx context.Context, title, slug string) error

	// ValidateUpdate 校验更新参数（排除自身 slug 唯一性）
	ValidateUpdate(ctx context.Context, id uint64, title, slug string) error

	// BuildArticle 构建新文章实体
	BuildArticle(authorID uint64, req BuildArticleReq) (*Article, error)

	// EnsureSlug 若 slug 为空则由 title 自动生成唯一 slug
	EnsureSlug(ctx context.Context, title, slug string, excludeID uint64) (string, error)
}

// BuildArticleReq 构建文章所需参数
type BuildArticleReq struct {
	Title        string
	Slug         string
	Summary      string
	Content      string
	ContentHTML  string
	Cover        string
	CategoryID   *uint64
	TagIDs       []uint64
	Status       Status
	IsTop        bool
	IsFeatured   bool
	AllowComment bool
}

// DomainService 文章领域服务实现
type DomainService struct {
	repo IArticleRepository
}

// NewDomainService 创建文章领域服务
func NewDomainService(repo IArticleRepository) IDomainService {
	return &DomainService{repo: repo}
}

func (s *DomainService) ValidateCreate(ctx context.Context, title, slug string) error {
	if title == "" {
		return apperrors.ErrInvalidParams("文章标题不能为空")
	}
	if slug != "" {
		if !utils.IsValidSlug(slug) {
			return apperrors.ErrInvalidParams("Slug 格式无效，只能包含小写字母、数字和连字符")
		}
		exists, err := s.repo.ExistsBySlug(ctx, slug, 0)
		if err != nil {
			return apperrors.ErrInternalError(err)
		}
		if exists {
			return apperrors.FromCode(apperrors.CodeSlugAlreadyExists)
		}
	}
	return nil
}

func (s *DomainService) ValidateUpdate(ctx context.Context, id uint64, title, slug string) error {
	if title == "" {
		return apperrors.ErrInvalidParams("文章标题不能为空")
	}
	if slug != "" {
		if !utils.IsValidSlug(slug) {
			return apperrors.ErrInvalidParams("Slug 格式无效，只能包含小写字母、数字和连字符")
		}
		exists, err := s.repo.ExistsBySlug(ctx, slug, id)
		if err != nil {
			return apperrors.ErrInternalError(err)
		}
		if exists {
			return apperrors.FromCode(apperrors.CodeSlugAlreadyExists)
		}
	}
	return nil
}

func (s *DomainService) BuildArticle(authorID uint64, req BuildArticleReq) (*Article, error) {
	if req.Title == "" {
		return nil, apperrors.ErrInvalidParams("文章标题不能为空")
	}

	if !req.Status.IsValid() {
		req.Status = StatusDraft
	}

	// 计算摘要和字数
	summary := AutoSummary(req.Content, req.Summary)
	wordCount := CountWords(req.Content)

	a := &Article{
		UUID:         uuid.NewString(),
		AuthorID:     authorID,
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		Slug:         req.Slug,
		Summary:      summary,
		Content:      req.Content,
		ContentHTML:  req.ContentHTML,
		Cover:        req.Cover,
		Status:       req.Status,
		IsTop:        req.IsTop,
		IsFeatured:   req.IsFeatured,
		AllowComment: req.AllowComment,
		WordCount:    wordCount,
		TagIDs:       req.TagIDs,
	}

	// 若初始状态为已发布，记录发布时间
	if a.Status == StatusPublished {
		now := time.Now()
		a.PublishedAt = &now
	}

	return a, nil
}

func (s *DomainService) EnsureSlug(ctx context.Context, title, slug string, excludeID uint64) (string, error) {
	if slug != "" {
		return slug, nil
	}
	base := utils.Slugify(title)
	if base == "" {
		base = "article"
	}
	candidate := base
	for i := 1; i <= 20; i++ {
		exists, err := s.repo.ExistsBySlug(ctx, candidate, excludeID)
		if err != nil {
			return "", apperrors.ErrInternalError(err)
		}
		if !exists {
			return candidate, nil
		}
		candidate = utils.SlugifyWithSuffix(title, i)
	}
	return "", apperrors.ErrInvalidParams("无法生成唯一 Slug，请手动指定")
}
