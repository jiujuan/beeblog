// Package validator 封装 go-playground/validator，
// 提供结构体验证、自定义规则注册和友好的中文错误消息。
package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// ─── 接口定义 ────────────────────────────────────────────────────────────────

// IValidator 验证器接口
type IValidator interface {
	// Validate 验证结构体，返回第一条错误（中文）
	Validate(s any) error
	// ValidateAll 验证结构体，返回所有错误
	ValidateAll(s any) []string
	// RegisterRule 注册自定义验证规则
	RegisterRule(tag string, fn validator.Func, msg string) error
}

// ─── 实现 ────────────────────────────────────────────────────────────────────

// Validator 封装 go-playground/validator
type Validator struct {
	validate *validator.Validate
	trans    ut.Translator
	once     sync.Once
}

var (
	instance *Validator
	mu       sync.Mutex
)

// Default 返回全局单例 Validator（线程安全）
func Default() *Validator {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = mustNew()
	}
	return instance
}

// New 创建新的 Validator 实例
func New() (*Validator, error) {
	v := &Validator{}
	if err := v.init(); err != nil {
		return nil, err
	}
	return v, nil
}

func mustNew() *Validator {
	v, err := New()
	if err != nil {
		panic(err)
	}
	return v
}

func (v *Validator) init() error {
	zhLocale := zh.New()
	uni := ut.New(zhLocale, zhLocale)

	trans, ok := uni.GetTranslator("zh")
	if !ok {
		return fmt.Errorf("validator: zh translator not found")
	}
	v.trans = trans

	validate := validator.New()

	// 使用 json tag 作为字段名（对前端更友好）
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	// 注册中文翻译
	if err := zhTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		return fmt.Errorf("validator: register zh translations: %w", err)
	}

	v.validate = validate
	return nil
}

// Validate 验证结构体，返回第一条中文错误信息
func (v *Validator) Validate(s any) error {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if ok := isValidationErrors(err, &errs); ok && len(errs) > 0 {
		return fmt.Errorf("%s", errs[0].Translate(v.trans))
	}
	return err
}

// ValidateAll 验证结构体，返回所有中文错误信息
func (v *Validator) ValidateAll(s any) []string {
	err := v.validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if ok := isValidationErrors(err, &errs); ok {
		msgs := make([]string, 0, len(errs))
		for _, e := range errs {
			msgs = append(msgs, e.Translate(v.trans))
		}
		return msgs
	}
	return []string{err.Error()}
}

// ValidateVar 验证单个变量
func (v *Validator) ValidateVar(field any, tag string) error {
	return v.validate.Var(field, tag)
}

// RegisterRule 注册自定义验证规则
// msg 示例："{0}必须是有效的手机号"
func (v *Validator) RegisterRule(tag string, fn validator.Func, msg string) error {
	if err := v.validate.RegisterValidation(tag, fn); err != nil {
		return fmt.Errorf("validator: register rule %q: %w", tag, err)
	}

	// 注册翻译
	if msg != "" {
		_ = v.validate.RegisterTranslation(tag, v.trans,
			func(ut ut.Translator) error {
				return ut.Add(tag, msg, true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(tag, fe.Field())
				return t
			},
		)
	}
	return nil
}

// ─── 辅助函数 ────────────────────────────────────────────────────────────────

func isValidationErrors(err error, target *validator.ValidationErrors) bool {
	if e, ok := err.(validator.ValidationErrors); ok {
		*target = e
		return true
	}
	return false
}

// Validate 全局便捷函数
func Validate(s any) error {
	return Default().Validate(s)
}

// ValidateAll 全局便捷函数，返回所有错误
func ValidateAll(s any) []string {
	return Default().ValidateAll(s)
}
