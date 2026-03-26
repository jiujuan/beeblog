// Package errors 提供统一的业务错误类型，支持错误码、HTTP 状态码与链式包装。
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError 业务错误，携带错误码和可选的底层原因
type AppError struct {
	Code    int    // 业务错误码
	Message string // 面向用户的错误描述
	Err     error  // 原始错误（内部使用，不对外暴露）
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Is / errors.As 链式解包
func (e *AppError) Unwrap() error {
	return e.Err
}

// ─── 构造函数 ────────────────────────────────────────────────────────────────

// New 使用错误码和自定义消息构建 AppError
func New(code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

// NewWithErr 使用错误码、消息和底层错误构建 AppError
func NewWithErr(code int, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

// FromCode 使用错误码（从注册表中查找默认消息）
func FromCode(code int) *AppError {
	return &AppError{Code: code, Message: Message(code)}
}

// Wrap 在已有错误上附加错误码（若 err 本身是 AppError 则直接返回）
func Wrap(code int, err error) *AppError {
	if err == nil {
		return nil
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return &AppError{Code: code, Message: Message(code), Err: err}
}

// ─── 常用快捷构造 ────────────────────────────────────────────────────────────

func ErrInternalError(err error) *AppError {
	return NewWithErr(CodeInternalError, Message(CodeInternalError), err)
}

func ErrInvalidParams(msg string) *AppError {
	return New(CodeInvalidParams, msg)
}

func ErrUnauthorized() *AppError {
	return FromCode(CodeUnauthorized)
}

func ErrForbidden() *AppError {
	return FromCode(CodeForbidden)
}

func ErrNotFound(msg string) *AppError {
	return New(CodeNotFound, msg)
}

func ErrUserNotFound() *AppError      { return FromCode(CodeUserNotFound) }
func ErrUserDisabled() *AppError      { return FromCode(CodeUserDisabled) }
func ErrPasswordIncorrect() *AppError { return FromCode(CodePasswordIncorrect) }
func ErrTokenExpired() *AppError      { return FromCode(CodeTokenExpired) }
func ErrTokenInvalid() *AppError      { return FromCode(CodeTokenInvalid) }
func ErrArticleNotFound() *AppError   { return FromCode(CodeArticleNotFound) }
func ErrCategoryNotFound() *AppError  { return FromCode(CodeCategoryNotFound) }
func ErrTagNotFound() *AppError       { return FromCode(CodeTagNotFound) }
func ErrCommentNotFound() *AppError   { return FromCode(CodeCommentNotFound) }
func ErrMediaNotFound() *AppError     { return FromCode(CodeMediaNotFound) }

// ─── 辅助函数 ────────────────────────────────────────────────────────────────

// IsAppError 判断 err 是否为 AppError
func IsAppError(err error) bool {
	var ae *AppError
	return errors.As(err, &ae)
}

// GetAppError 从 err 中提取 AppError，若不是则返回 nil
func GetAppError(err error) *AppError {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return nil
}

// HTTPStatus 将业务错误码映射到 HTTP 状态码
func HTTPStatus(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeInvalidParams:
		return http.StatusBadRequest
	case CodeUnauthorized, CodeTokenExpired, CodeTokenInvalid, CodeRefreshTokenExpired:
		return http.StatusUnauthorized
	case CodeForbidden, CodeUserDisabled:
		return http.StatusForbidden
	case CodeUserAlreadyExists, CodeEmailAlreadyExists:
		return http.StatusConflict
	case CodePasswordIncorrect:
		return http.StatusUnauthorized
	case CodeNotFound,
		CodeUserNotFound, CodeArticleNotFound,
		CodeCategoryNotFound, CodeTagNotFound,
		CodeCommentNotFound, CodeMediaNotFound:
		return http.StatusNotFound
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
