package validations

import (
	"payment-service/pkg/errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationError {
	var validationErrors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: getErrorMessage(err),
			})
		}
	}

	return validationErrors
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters"
	case "gt":
		return err.Field() + " must be greater than " + err.Param()
	case "gte":
		return err.Field() + " must be greater than or equal to " + err.Param()
	case "lt":
		return err.Field() + " must be less than " + err.Param()
	case "lte":
		return err.Field() + " must be less than or equal to " + err.Param()
	case "uuid":
		return err.Field() + " must be a valid UUID"
	default:
		return err.Field() + " is invalid"
	}
}

func ValidatePaymentAmount(amount float64) error {
	if amount <= 0 {
		return errors.ErrInvalidAmount.WithDetails("Amount must be greater than 0")
	}
	return nil
}

func ValidateCurrency(currency string) error {
	if currency != "IDR" {
		return errors.ErrBadRequest.WithDetails("Only IDR currency is supported")
	}
	return nil
}