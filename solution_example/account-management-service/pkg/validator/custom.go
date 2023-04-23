package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
)

const (
	passwordMinLength = 8
	passwordMaxLength = 32
	passwordMinLower  = 1
	passwordMinUpper  = 1
	passwordMinDigit  = 1
	passwordMinSymbol = 1
)

var (
	lengthRegexp    = regexp.MustCompile(fmt.Sprintf(`^.{%d,%d}$`, passwordMinLength, passwordMaxLength))
	lowerCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[a-z]{%d,}`, passwordMinLower))
	upperCaseRegexp = regexp.MustCompile(fmt.Sprintf(`[A-Z]{%d,}`, passwordMinUpper))
	digitRegexp     = regexp.MustCompile(fmt.Sprintf(`[0-9]{%d,}`, passwordMinDigit))
	symbolRegexp    = regexp.MustCompile(fmt.Sprintf(`[!@#$%%^&*]{%d,}`, passwordMinSymbol))
)

type CustomValidator struct {
	v         *validator.Validate
	passwdErr error
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{v: v}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.RegisterValidation("password", cv.passwordValidate)
	if err != nil {
		panic(err)
	}

	return cv
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.v.Struct(i)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]

		return cv.newValidationError(fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), fieldErr.Param())
	}
	return nil
}

func (cv *CustomValidator) newValidationError(field string, value interface{}, tag string, param string) error {
	switch tag {
	case "required":
		return fmt.Errorf("field %s is required", field)
	case "email":
		return fmt.Errorf("field %s must be a valid email address", field)
	case "password":
		return cv.passwdErr
	case "min":
		return fmt.Errorf("field %s must be at least %s characters", field, param)
	case "max":
		return fmt.Errorf("field %s must be at most %s characters", field, param)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}

func (cv *CustomValidator) passwordValidate(fl validator.FieldLevel) bool {
	// check if the field is a string
	if fl.Field().Kind() != reflect.String {
		cv.passwdErr = fmt.Errorf("field %s must be a string", fl.FieldName())
		return false
	}

	// get the value of the field
	fieldValue := fl.Field().String()

	// check regexp matching
	if ok := lengthRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must be between %d and %d characters", fl.FieldName(), passwordMinLength, passwordMaxLength)
		return false
	} else if ok = lowerCaseRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must contain at least %d lowercase letter(s)", fl.FieldName(), passwordMinLower)
		return false
	} else if ok = upperCaseRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must contain at least %d uppercase letter(s)", fl.FieldName(), passwordMinUpper)
		return false
	} else if ok = digitRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must contain at least %d digit(s)", fl.FieldName(), passwordMinDigit)
		return false
	} else if ok = symbolRegexp.MatchString(fieldValue); !ok {
		cv.passwdErr = fmt.Errorf("field %s must contain at least %d special character(s)", fl.FieldName(), passwordMinSymbol)
		return false
	}

	return true
}
