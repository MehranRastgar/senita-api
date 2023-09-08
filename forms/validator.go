package forms

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// DefaultValidator ...
type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var validate = validator.New()

// NewValidator creates a new instance of the DefaultValidator.
func NewValidator() *DefaultValidator {
	return &DefaultValidator{}
}

// RegisterCustomValidations registers custom validation rules.
func (v *DefaultValidator) RegisterCustomValidations() {
	v.lazyinit()

	// Register custom validation rules here
	v.validate.RegisterValidation("fullName", ValidateFullName)
}

// Validate validates a struct using the validator.
func (v *DefaultValidator) Validate(ctx *fiber.Ctx, i interface{}) error {
	if kindOfData(i) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(i); err != nil {
			return err
		}
	}
	return nil
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// Register custom validation rules here
		v.validate.RegisterValidation("fullName", ValidateFullName)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// ValidateFullName implements validator.Func
func ValidateFullName(fl validator.FieldLevel) bool {
	// Remove the extra space
	space := regexp.MustCompile(`\s+`)
	name := space.ReplaceAllString(fl.Field().String(), " ")

	// Remove trailing spaces
	name = strings.TrimSpace(name)

	// To support all possible languages
	matched, _ := regexp.Match(`^[^±!@£$%^&*_+§¡€#¢§¶•ªº«\\/<>?:;'"|=.,0123456789]{3,20}$`, []byte(name))
	return matched
}
