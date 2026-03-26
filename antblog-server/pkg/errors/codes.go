package errors

// 业务错误码常量定义
// 规范：
//   1xxxx - 系统级错误
//   2xxxx - 用户模块
//   3xxxx - 文章模块
//   4xxxx - 分类/标签模块
//   5xxxx - 评论模块
//   6xxxx - 媒体模块

const (
	// ── 系统级 ──────────────────────────────────────────
	CodeSuccess         = 0
	CodeInternalError   = 10000
	CodeInvalidParams   = 10001
	CodeUnauthorized    = 10002
	CodeForbidden       = 10003
	CodeNotFound        = 10004
	CodeTooManyRequests = 10005
	CodeServiceUnavail  = 10006

	// ── 用户模块 ─────────────────────────────────────────
	CodeUserNotFound       = 20001
	CodeUserAlreadyExists  = 20002
	CodePasswordIncorrect  = 20003
	CodeUserDisabled       = 20004
	CodeTokenExpired       = 20005
	CodeTokenInvalid       = 20006
	CodeRefreshTokenExpired = 20007
	CodeEmailAlreadyExists = 20008

	// ── 文章模块 ─────────────────────────────────────────
	CodeArticleNotFound    = 30001
	CodeArticleDeleted     = 30002
	CodeArticleNotPublished = 30003

	// ── 分类/标签模块 ────────────────────────────────────
	CodeCategoryNotFound  = 40001
	CodeCategoryHasArticle = 40002
	CodeTagNotFound       = 40003
	CodeSlugAlreadyExists = 40004

	// ── 评论模块 ─────────────────────────────────────────
	CodeCommentNotFound   = 50001
	CodeCommentNotAllowed = 50002

	// ── 媒体模块 ─────────────────────────────────────────
	CodeMediaNotFound     = 60001
	CodeMediaUploadFailed = 60002
	CodeMediaInvalidType  = 60003
	CodeMediaTooLarge     = 60004
)

// codeMessages 错误码对应的默认消息
var codeMessages = map[int]string{
	CodeSuccess:         "success",
	CodeInternalError:   "内部服务器错误",
	CodeInvalidParams:   "请求参数错误",
	CodeUnauthorized:    "未授权，请先登录",
	CodeForbidden:       "无访问权限",
	CodeNotFound:        "资源不存在",
	CodeTooManyRequests: "请求过于频繁，请稍后再试",
	CodeServiceUnavail:  "服务暂不可用",

	CodeUserNotFound:        "用户不存在",
	CodeUserAlreadyExists:   "用户名已存在",
	CodePasswordIncorrect:   "密码错误",
	CodeUserDisabled:        "账号已被禁用",
	CodeTokenExpired:        "Token 已过期",
	CodeTokenInvalid:        "Token 无效",
	CodeRefreshTokenExpired: "刷新 Token 已过期，请重新登录",
	CodeEmailAlreadyExists:  "邮箱已被注册",

	CodeArticleNotFound:     "文章不存在",
	CodeArticleDeleted:      "文章已被删除",
	CodeArticleNotPublished: "文章未发布",

	CodeCategoryNotFound:   "分类不存在",
	CodeCategoryHasArticle: "分类下存在文章，无法删除",
	CodeTagNotFound:        "标签不存在",
	CodeSlugAlreadyExists:  "Slug 已存在",

	CodeCommentNotFound:   "评论不存在",
	CodeCommentNotAllowed: "该文章不允许评论",

	CodeMediaNotFound:     "媒体资源不存在",
	CodeMediaUploadFailed: "文件上传失败",
	CodeMediaInvalidType:  "不支持的文件类型",
	CodeMediaTooLarge:     "文件大小超出限制",
}

// Message 根据错误码返回默认消息
func Message(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
