package controller

import (
	"backend/models"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// Define a global translator
var trans ut.Translator

// InitTrans initializes the translator
func InitTrans(locale string) (err error) {
	// Modify the Gin framework's Validator engine property for customization
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register a custom method to get the JSON tag
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		// Register custom validation method for SignUpParam
		v.RegisterStructValidation(SignUpParamStructLevelValidation, models.RegisterForm{})
		zhT := zh.New() // Chinese translator
		enT := en.New() // English translator
		// The first parameter is the fallback locale
		// The subsequent parameters are the locales to support (multiple supported)
		// uni := ut.New(zhT, zhT) would also work
		uni := ut.New(enT, zhT, enT)
		// locale is usually determined by the 'Accept-Language' header in the HTTP request
		var ok bool
		// You can also use uni.FindTranslator(...) to search for multiple locales
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		// Register the translator
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// Define a custom method to remove the struct name prefix:
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// SignUpParamStructLevelValidation custom SignUpParam structure validation function
func SignUpParamStructLevelValidation(sl validator.StructLevel) {
	su := sl.Current().Interface().(models.RegisterForm)

	if su.Password != su.ConfirmPassword {
		// Output the error message, the last parameter is the passed param
		sl.ReportError(su.ConfirmPassword, "confirm_password", "ConfirmPassword", "eqfield", "password")
	}
}
