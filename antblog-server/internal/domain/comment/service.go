package comment

import (
	"regexp"
	"strings"
	"unicode/utf8"

	apperrors "antblog/pkg/errors"
)

const (
	maxContentLen  = 2000 // 评论字数上限
	minContentLen  = 1
	maxNicknameLen = 64
	maxEmailLen    = 128
)

var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// IDomainService 评论领域服务接口
type IDomainService interface {
	// ValidateContent 校验评论内容
	ValidateContent(content string) error

	// ValidateGuest 校验游客信息（昵称/邮箱）
	ValidateGuest(nickname, email string) error

	// BuildComment 构建评论实体
	BuildComment(req BuildCommentReq) (*Comment, error)
}

// BuildCommentReq 构建评论所需参数
type BuildCommentReq struct {
	ArticleID uint64
	UserID    *uint64
	ParentID  *uint64
	Parent    *Comment // 父评论实体，用于计算 RootID
	ReplyToID *uint64
	Nickname  string
	Email     string
	Content   string
	IP        string
	UserAgent string
}

// DomainService 评论领域服务实现
type DomainService struct{}

// NewDomainService 创建评论领域服务
func NewDomainService() IDomainService {
	return &DomainService{}
}

func (s *DomainService) ValidateContent(content string) error {
	content = strings.TrimSpace(content)
	if utf8.RuneCountInString(content) < minContentLen {
		return apperrors.ErrInvalidParams("评论内容不能为空")
	}
	if utf8.RuneCountInString(content) > maxContentLen {
		return apperrors.ErrInvalidParams("评论内容不得超过 2000 个字符")
	}
	return nil
}

func (s *DomainService) ValidateGuest(nickname, email string) error {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return apperrors.ErrInvalidParams("游客昵称不能为空")
	}
	if utf8.RuneCountInString(nickname) > maxNicknameLen {
		return apperrors.ErrInvalidParams("昵称不得超过 64 个字符")
	}
	if email != "" {
		if utf8.RuneCountInString(email) > maxEmailLen || !emailRegexp.MatchString(email) {
			return apperrors.ErrInvalidParams("邮箱格式无效")
		}
	}
	return nil
}

func (s *DomainService) BuildComment(req BuildCommentReq) (*Comment, error) {
	if err := s.ValidateContent(req.Content); err != nil {
		return nil, err
	}

	// 游客评论必须提供昵称
	if req.UserID == nil {
		if err := s.ValidateGuest(req.Nickname, req.Email); err != nil {
			return nil, err
		}
	}

	c := &Comment{
		ArticleID: req.ArticleID,
		UserID:    req.UserID,
		ParentID:  req.ParentID,
		ReplyToID: req.ReplyToID,
		Nickname:  strings.TrimSpace(req.Nickname),
		Email:     strings.TrimSpace(req.Email),
		Content:   strings.TrimSpace(req.Content),
		IP:        req.IP,
		UserAgent: req.UserAgent,
		Status:    StatusPending, // 默认待审核
	}

	// 由父评论推导 RootID
	if req.Parent != nil {
		c.SetRootID(req.Parent)
	}

	return c, nil
}
