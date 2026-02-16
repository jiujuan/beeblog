package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 工具库: 提供常用工具函数，如响应格式、分页、字符串处理等。

// ========== 响应结构 ==========
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(msg string, data interface{}) Response {
	return Response{
		Code:    200,
		Message: msg,
		Data:    data,
	}
}

// Fail 失败响应
func Fail(code int, msg string) Response {
	return Response{
		Code:    code,
		Message: msg,
	}
}

// ========== 分页 ==========

type Pagination struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

func (p *Pagination) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // 限制最大值，防止性能问题
	}
	return p.PageSize
}

// ========== 字符串工具 ==========

// RandomString 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// UUID 生成 UUID
func UUID() string {
	return uuid.New().String()
}

// IsEmail 验证邮箱格式
func IsEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// Truncate 截断字符串
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// ========== 时间工具 ==========

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseTime 解析时间
func ParseTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}

// ========== Slice 工具 ==========

// Contains 检查 slice 是否包含元素
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveEmpty 移除空字符串
func RemoveEmpty(slice []string) []string {
	var result []string
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}

// ========== 错误处理 ==========

// WrapError 包装错误信息
func WrapError(msg string, err error) error {
	if err == nil {
		return fmt.Errorf("%s", msg)
	}
	return fmt.Errorf("%s: %w", msg, err)
}
