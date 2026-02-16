package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"beeblog/pkg/utils"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

// Init 初始化验证器
func Init(locale string) error {
	// 获取 gin 默认的 validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
	} else {
		validate = validator.New()
	}

	// 注册自定义验证规则
	registerCustomValidators()

	// 设置翻译器
	if err := registerTranslator(locale); err != nil {
		return err
	}

	return nil
}

// Validate 结构体验证
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// 翻译错误信息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	// 只返回第一个错误
	if len(errs) > 0 {
		return fmt.Errorf("%s", errs[0].Translate(trans))
	}

	return nil
}

// ValidateVar 变量验证
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

func registerTranslator(locale string) error {
	uni := ut.New(en.New(), zh.New())
	var ok bool
	trans, ok = uni.GetTranslator(locale)
	if !ok {
		trans = uni.GetFallback()
	}

	var err error
	switch locale {
	case "zh":
		err = zh_translations.RegisterDefaultTranslations(validate, trans)
	default:
		err = en_translations.RegisterDefaultTranslations(validate, trans)
	}

	if err != nil {
		return err
	}

	// 注册自定义翻译
	_ = validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0}为必填字段", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	return nil
}

func registerCustomValidators() {
	// 注册自定义验证规则：检查密码强度
	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 6 || len(password) > 20 {
			return false
		}
		// 至少包含一个字母和一个数字
		matched, _ := regexp.MatchString("^(?=.*[A-Za-z])(?=.*\\d).+$", password)
		return matched
	})
}

// GetValidMsg 获取结构体中的 label 标签作为字段名
func GetValidMsg(err error, obj interface{}) string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	getField := func(fieldName string) string {
		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return fieldName
		}

		field, found := t.FieldByName(fieldName)
		if !found {
			return fieldName
		}

		label := field.Tag.Get("label")
		if label != "" {
			return label
		}
		return fieldName
	}

	if len(errs) > 0 {
		fe := errs[0]
		return fmt.Sprintf("%s %s", getField(fe.Field()), getErrorMsg(fe))
	}

	return err.Error()
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "不能为空"
	case "email":
		return "格式不正确"
	case "password":
		return "必须包含字母和数字，长度6-20位"
	default:
		return fe.Error()
	}
}

// 封装给 utils 使用
func Translate(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok && len(errs) > 0 {
		return errs[0].Translate(trans)
	}
	return err.Error()
}
