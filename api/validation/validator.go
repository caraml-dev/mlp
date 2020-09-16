package validation

import (
	"regexp"
	"reflect"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-playground/validator"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	
)

var (
	english = en_US.New()
	trans   = ut.New(english, english)
	EN, _   = trans.GetTranslator("en")
	
	subdomainRffc1123 = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-]+[a-zA-Z0-9]$`)
)

func NewValidator() *validator.Validate {
	instance := validator.New()
	instance.RegisterValidation("notblank", NotBlank)
	instance.RegisterValidation("subdomain_rfc1123", isRFC1123SubDomain)
	instance.RegisterValidation("url", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		switch field.Kind() {
		case reflect.String:
			return govalidator.IsURL(field.String())
		default:
			return false
		}
	})

	instance.RegisterTranslation("required", EN, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			fld := fe.StructField()
			t, err := ut.T(fe.Tag(), fld)
			if err != nil {
				return fe.(error).Error()
			}
			return t
		})

	instance.RegisterTranslation("min", EN, func(ut ut.Translator) error {
		return ut.Add("min", "{0} should be more than {1} characters", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			fld := fe.StructField()
			param := fe.Param()
			t, err := ut.T(fe.Tag(), fld, param)
			if err != nil {
				return fe.(error).Error()
			}
			return t
		})

	instance.RegisterTranslation("max", EN, func(ut ut.Translator) error {
		return ut.Add("max", "{0} should be less than {1} characters", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			fld := fe.StructField()
			param := fe.Param()
			t, err := ut.T(fe.Tag(), fld, param)
			if err != nil {
				return fe.(error).Error()
			}
			return t
		})

	instance.RegisterTranslation("subdomain_rfc1123", EN, func(ut ut.Translator) error {
		return ut.Add("subdomain_rfc1123", "{0} should be a valid RFC1123 sub-domain", true)
	},
		func(ut ut.Translator, fe validator.FieldError) string {
			fld := fe.StructField()
			t, err := ut.T(fe.Tag(), fld)
			if err != nil {
				return fe.(error).Error()
			}
			return t
		})
		
	return instance
}

func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isRFC1123SubDomain(fl validator.FieldLevel) bool {
	return subdomainRffc1123.MatchString(fl.Field().String())
}