package apierrors

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ResponseError struct {
	Message string `json:"message"`
}

func New(errString string) ResponseError {
	return ResponseError{
		Message: errString,
	}
}

type CustomValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func (cv *CustomValidator) Validate(i interface{}) error {

	errorStrings := make([]string, 0)

	err := cv.validator.Struct(i)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errorStrings = append(errorStrings, e.Translate(cv.trans))
		}
		return errors.New(strings.Join(errorStrings, "; "))
	}

	return nil
}

func NewCustomValidator() *CustomValidator {
	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validator: v, trans: trans}
}
