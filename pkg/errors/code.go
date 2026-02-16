package errors

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

// 定义常用错误码
var (
	Success         = New(0, "success")
	ServerInternal  = New(500, "服务器内部错误")
	InvalidParams   = New(400, "参数错误")
	Unauthorized    = New(401, "未授权")
	Forbidden       = New(403, "禁止访问")
	NotFound        = New(404, "资源不存在")
	TooManyRequests = New(429, "请求过于频繁")

	// 用户相关错误 (1000-1999)
	UserNotFound      = New(1001, "用户不存在")
	UserAlreadyExists = New(1002, "用户已存在")
	InvalidPassword   = New(1003, "密码错误")
	InvalidToken      = New(1004, "Token 无效")
	TokenExpired      = New(1005, "Token 已过期")

	// 文章相关错误 (2000-2999)
	PostNotFound     = New(2001, "文章不存在")
	CategoryNotFound = New(2002, "分类不存在")
	TagNotFound      = New(2003, "标签不存在")
)
