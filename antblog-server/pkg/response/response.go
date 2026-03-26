// Package response 提供统一的 HTTP 响应格式封装。
// 响应结构：{ "code": 0, "msg": "success", "data": {...} }
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "antblog/pkg/errors"
)

// Response 统一响应体
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// PageData 分页响应数据
type PageData[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// NewPageData 构建分页数据
func NewPageData[T any](list []T, total int64, page, pageSize int) *PageData[T] {
	if list == nil {
		list = make([]T, 0)
	}
	return &PageData[T]{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

// ─── 成功响应 ────────────────────────────────────────────────────────────────

// OK 返回成功响应，携带数据
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code: apperrors.CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// OKWithMsg 返回成功响应，携带自定义消息和数据
func OKWithMsg(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusOK, Response{
		Code: apperrors.CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// Created 返回 201 创建成功响应
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Code: apperrors.CodeSuccess,
		Msg:  "created",
		Data: data,
	})
}

// NoContent 返回 204 无内容响应（删除操作）
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ─── 失败响应 ────────────────────────────────────────────────────────────────

// Fail 返回业务失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(apperrors.HTTPStatus(code), Response{
		Code: code,
		Msg:  msg,
	})
}

// FailWithError 从 AppError 自动提取码和消息
func FailWithError(c *gin.Context, err error) {
	ae := apperrors.GetAppError(err)
	if ae != nil {
		c.JSON(apperrors.HTTPStatus(ae.Code), Response{
			Code: ae.Code,
			Msg:  ae.Message,
		})
		return
	}
	// 非 AppError，视为内部错误
	c.JSON(http.StatusInternalServerError, Response{
		Code: apperrors.CodeInternalError,
		Msg:  apperrors.Message(apperrors.CodeInternalError),
	})
}

// BadRequest 400 参数错误
func BadRequest(c *gin.Context, msg string) {
	Fail(c, apperrors.CodeInvalidParams, msg)
}

// Unauthorized 401 未授权
func Unauthorized(c *gin.Context) {
	Fail(c, apperrors.CodeUnauthorized, apperrors.Message(apperrors.CodeUnauthorized))
}

// Forbidden 403 禁止访问
func Forbidden(c *gin.Context) {
	Fail(c, apperrors.CodeForbidden, apperrors.Message(apperrors.CodeForbidden))
}

// NotFound 404 资源不存在
func NotFound(c *gin.Context, msg string) {
	Fail(c, apperrors.CodeNotFound, msg)
}

// InternalError 500 内部错误
func InternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: apperrors.CodeInternalError,
		Msg:  apperrors.Message(apperrors.CodeInternalError),
	})
}
