package validations

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance with custom configuration
func New() *Validator {
	v := validator.New()

	// Register custom tag name function for JSON field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}
		return name
	})

	// Register custom validators
	registerCustomValidators(v)

	return &Validator{
		validate: v,
	}
}

// Engine returns the underlying validator engine
func (v *Validator) Engine() interface{} {
	return v.validate
}

// Reset removes all validations and registered custom validators
// Use with caution - this is primarily for testing
func (v *Validator) Reset() {
	v.validate = validator.New()
	registerCustomValidators(v.validate)
}
