package user

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ─── Email 值对象 ────────────────────────────────────────────────────────────

var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Email 邮箱值对象（不可变）
type Email struct {
	value string
}

// NewEmail 构造并校验 Email
func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return Email{}, errors.New("邮箱不能为空")
	}
	if len(raw) > 128 {
		return Email{}, errors.New("邮箱长度不能超过128位")
	}
	if !emailRegexp.MatchString(raw) {
		return Email{}, errors.New("邮箱格式不正确")
	}
	return Email{value: raw}, nil
}

// String 返回邮箱字符串
func (e Email) String() string { return e.value }

// Equals 判断是否相等
func (e Email) Equals(other Email) bool { return e.value == other.value }

// ─── Password 值对象 ─────────────────────────────────────────────────────────

const (
	PasswordMinLen = 6
	PasswordMaxLen = 72 // bcrypt 最大有效长度
)

// Password 密码值对象，内部存储明文（仅在注册/修改密码时短暂存在）
type Password struct {
	raw string
}

// NewPassword 构造并校验密码强度
func NewPassword(raw string) (Password, error) {
	if raw == "" {
		return Password{}, errors.New("密码不能为空")
	}
	length := utf8.RuneCountInString(raw)
	if length < PasswordMinLen {
		return Password{}, errors.New("密码长度不能少于6位")
	}
	if length > PasswordMaxLen {
		return Password{}, errors.New("密码长度不能超过72位")
	}
	return Password{raw: raw}, nil
}

// Raw 返回原始密码（仅供加密使用，禁止日志输出）
func (p Password) Raw() string { return p.raw }

// ─── Username 值对象 ─────────────────────────────────────────────────────────

var usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_\-]{3,32}$`)

// Username 用户名值对象
type Username struct {
	value string
}

// NewUsername 构造并校验用户名
func NewUsername(raw string) (Username, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Username{}, errors.New("用户名不能为空")
	}
	if !usernameRegexp.MatchString(raw) {
		return Username{}, errors.New("用户名只能包含字母、数字、下划线和连字符，长度3-32位")
	}
	return Username{value: raw}, nil
}

// String 返回用户名字符串
func (u Username) String() string { return u.value }
