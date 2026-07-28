package validations

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// ValidationErrors represents a map of field names to error messages
type ValidationErrors map[string][]string

// Error implements the error interface for ValidationErrors
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation failed"
	}
	var errs []string
	for field, messages := range ve {
		for _, msg := range messages {
			errs = append(errs, fmt.Sprintf("%s: %s", field, msg))
		}
	}
	return strings.Join(errs, "; ")
}

// HasErrors returns true if there are validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// GetFieldErrors returns errors for a specific field
func (ve ValidationErrors) GetFieldErrors(field string) []string {
	if errors, exists := ve[field]; exists {
		return errors
	}
	return nil
}

// FirstError returns the first error message
func (ve ValidationErrors) FirstError() string {
	for _, errors := range ve {
		if len(errors) > 0 {
			return errors[0]
		}
	}
	return ""
}

// ToMap converts ValidationErrors to a simple map[string][]string
func (ve ValidationErrors) ToMap() map[string][]string {
	return map[string][]string(ve)
}

// Merge combines two ValidationErrors
func (ve ValidationErrors) Merge(other ValidationErrors) ValidationErrors {
	for field, errors := range other {
		if _, exists := ve[field]; !exists {
			ve[field] = []string{}
		}
		ve[field] = append(ve[field], errors...)
	}
	return ve
}

// Validate validates a struct and returns ValidationErrors
func (v *Validator) Validate(s interface{}) ValidationErrors {
	err := v.validate.Struct(s)
	return parseValidationErrors(err)
}

// ValidateWithContext validates a struct with context
func (v *Validator) ValidateWithContext(ctx context.Context, s interface{}) ValidationErrors {
	err := v.validate.StructCtx(ctx, s)
	return parseValidationErrors(err)
}

// ValidateAndReturnError validates a struct and returns a simple error string
func (v *Validator) ValidateAndReturnError(s interface{}) error {
	errs := v.Validate(s)
	if errs.HasErrors() {
		return errs
	}
	return nil
}

// ValidateAndPanic validates a struct and panics if validation fails
// Use this in places where validation MUST pass
func (v *Validator) ValidateAndPanic(s interface{}) {
	err := v.validate.Struct(s)
	if err != nil {
		panic(err)
	}
}

// ==================== Error Parsing ====================

// ParseValidationErrors converts a validator error to ValidationErrors map
// This is the public version that can be used from controllers
func ParseValidationErrors(err error) ValidationErrors {
	return parseValidationErrors(err)
}

// parseValidationErrors converts validator error to ValidationErrors map
func parseValidationErrors(err error) ValidationErrors {
	errors := make(ValidationErrors)

	if err == nil {
		return errors
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		// If it's not a validation error, return as single error
		errors["_error"] = []string{err.Error()}
		return errors
	}

	for _, fe := range validationErrs {
		field := fe.Field()
		if _, exists := errors[field]; !exists {
			errors[field] = []string{}
		}
		errors[field] = append(errors[field], validationErrorMessage(fe))
	}

	return errors
}

// validationErrorMessage converts a field error to a human-readable message
func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	// Required validations
	case "required":
		return fmt.Sprintf("%s is required", toSentenceCase(fe.Field()))
	case "required_if":
		return fmt.Sprintf("%s is required", toSentenceCase(fe.Field()))
	case "required_unless":
		return fmt.Sprintf("%s is required", toSentenceCase(fe.Field()))
	case "required_with":
		return fmt.Sprintf("%s is required", toSentenceCase(fe.Field()))
	case "required_without":
		return fmt.Sprintf("%s is required", toSentenceCase(fe.Field()))

	// String validations
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters", toSentenceCase(fe.Field()), fe.Param())
		}
		return fmt.Sprintf("%s must be at least %s", toSentenceCase(fe.Field()), fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters", toSentenceCase(fe.Field()), fe.Param())
		}
		return fmt.Sprintf("%s must be at most %s", toSentenceCase(fe.Field()), fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", toSentenceCase(fe.Field()), fe.Param())

	// Email validation
	case "email":
		return "Invalid email address"

	// URL validation
	case "url":
		return "Invalid URL format"

	// URI validation
	case "uri":
		return "Invalid URI format"

	// Numeric validations
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", toSentenceCase(fe.Field()), fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", toSentenceCase(fe.Field()), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", toSentenceCase(fe.Field()), fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", toSentenceCase(fe.Field()), fe.Param())
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", toSentenceCase(fe.Field()), fe.Param())
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", toSentenceCase(fe.Field()), fe.Param())

	// Number range validations
	case "min_value":
		return fmt.Sprintf("%s must be at least %s", toSentenceCase(fe.Field()), fe.Param())
	case "max_value":
		return fmt.Sprintf("%s must be at most %s", toSentenceCase(fe.Field()), fe.Param())
	case "eq_value":
		return fmt.Sprintf("%s must be equal to %s", toSentenceCase(fe.Field()), fe.Param())

	// String format validations
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", toSentenceCase(fe.Field()))
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", toSentenceCase(fe.Field()))
	case "numeric":
		return fmt.Sprintf("%s must contain only numbers", toSentenceCase(fe.Field()))
	case "alphanumeric":
		return fmt.Sprintf("%s must contain only letters and numbers", toSentenceCase(fe.Field()))

	// UUID validation
	case "uuid":
		return "Invalid UUID format"
	case "uuid3":
		return "Invalid UUID3 format"
	case "uuid4":
		return "Invalid UUID4 format"
	case "uuid5":
		return "Invalid UUID5 format"

	// Date validation
	case "datetime":
		return "Invalid date/time format"

	// JSON validation
	case "json":
		return "Invalid JSON format"

	// File path validations
	case "filepath":
		return "Invalid file path"
	case "file":
		return "File is required or invalid"
	case "mime_type":
		return fmt.Sprintf("%s has an invalid file type", toSentenceCase(fe.Field()))

	// IP validations
	case "ip":
		return "Invalid IP address"
	case "ipv4":
		return "Invalid IPv4 address"
	case "ipv6":
		return "Invalid IPv6 address"
	case "cidr":
		return "Invalid CIDR notation"
	case "cidr_ipv4":
		return "Invalid CIDR notation (IPv4)"
	case "cidr_ipv6":
		return "Invalid CIDR notation (IPv6)"

	// Country codes
	case "iso3166_1_alpha2":
		return "Invalid country code (ISO 3166-1 alpha-2)"
	case "iso3166_1_alpha3":
		return "Invalid country code (ISO 3166-1 alpha-3)"
	case "iso3166_1_alpha4":
		return "Invalid country code (ISO 3166-1 alpha-4)"
	case "iso3166_1_numeric":
		return "Invalid country code (ISO 3166-1 numeric)"

	// Currency codes
	case "iso4217":
		return "Invalid currency code (ISO 4217)"

	// Base64 validation
	case "base64":
		return "Invalid base64 encoding"
	case "base64url":
		return "Invalid base64 URL encoding"

	// Credit card validation
	case "credit_card":
		return "Invalid credit card number"

	// DNS/RFC validation
	case "dns":
		return "Invalid DNS name"
	case "fqdn":
		return "Invalid fully qualified domain name"
	case "hostname":
		return "Invalid hostname"
	case "hostname_rfc1123":
		return "Invalid hostname (RFC 1123)"
	case "url_safe":
		return "Invalid URL-safe string"

	// Eth addr validation
	case "eth_addr":
		return "Invalid Ethereum address"

	// Hex color validation
	case "hexcolor":
		return "Invalid hex color code"
	case "hex":
		return "Invalid hexadecimal value"

	// Print ASCII validation
	case "printascii":
		return "Invalid printable ASCII characters"

	// Latitude/Longitude validations
	case "latitude":
		return "Invalid latitude value"
	case "longitude":
		return "Invalid longitude value"

	// Phone number validation
	case "e164":
		return "Invalid phone number (E.164 format)"
	case "phone":
		return "Invalid phone number"

	// Postal code validation
	case "postcode_iso3166_alpha2":
		return "Invalid postal code for the specified country"

	// Semantic version validation
	case "semver":
		return "Invalid semantic version"

	// Urn RFC 2141 validation
	case "urn_rfc2141":
		return "Invalid URN (RFC 2141)"

	// Password strength validations
	case "password":
		return "Password does not meet requirements"

	// String contains/excludes validations
	case "contains":
		return fmt.Sprintf("%s must contain %s", toSentenceCase(fe.Field()), fe.Param())
	case "excludes":
		return fmt.Sprintf("%s must not contain %s", toSentenceCase(fe.Field()), fe.Param())
	case "containsany":
		return fmt.Sprintf("%s must contain at least one of: %s", toSentenceCase(fe.Field()), fe.Param())
	case "excludesall":
		return fmt.Sprintf("%s must not contain any of: %s", toSentenceCase(fe.Field()), fe.Param())

	// Custom validation
	case "matches":
		return fmt.Sprintf("%s does not match the required format", toSentenceCase(fe.Field()))

	// Custom validators
	case "phone_id":
		return "Invalid Indonesian phone number"
	case "strong_password":
		return "Password must contain uppercase, lowercase, number, and special character"
	case "date_range":
		return "End date must be after start date"
	case "unique_items":
		return "All items must be unique"

	default:
		return fmt.Sprintf("%s is invalid", toSentenceCase(fe.Field()))
	}
}

// ==================== Helpers ====================

// toSentenceCase converts a field name to sentence case
// e.g., "FirstName" -> "first name", "EmailAddress" -> "email address"
func toSentenceCase(s string) string {
	// Handle common abbreviations
	replacer := strings.NewReplacer(
		"ID", "ID",
		"URL", "URL",
		"API", "API",
		"HTML", "HTML",
		"JSON", "JSON",
		"XML", "XML",
		"SMTP", "SMTP",
		"FTP", "FTP",
		"HTTP", "HTTP",
		"HTTPS", "HTTPS",
		"SQL", "SQL",
		"No", "Number",
	)

	s = replacer.Replace(s)

	// Add space before uppercase letters (except for acronyms)
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			// Check if next char is also uppercase (acronym)
			nextIsUpper := i+1 < len(s) && s[i+1] >= 'A' && s[i+1] <= 'Z'
			if !nextIsUpper {
				result.WriteRune(' ')
			}
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}

// isNumeric checks if a string contains only numeric characters
func isNumeric(s string) bool {
	for _, char := range s {
		if char < '0' || char > '9' {
			return false
		}
	}
	return len(s) > 0
}

// ==================== Custom Validators ====================

// registerCustomValidators registers all custom validators
func registerCustomValidators(v *validator.Validate) {
	_ = v.RegisterValidation("phone_id", validatePhoneIndonesia)
	_ = v.RegisterValidation("strong_password", validateStrongPassword)
	_ = v.RegisterValidation("date_range", validateDateRange)
	_ = v.RegisterValidation("unique_items", validateUniqueItems)
}

// validatePhoneIndonesia validates Indonesian phone numbers
// Accepts: +62 xxx, 08xx xxx, 08xxxxxxxxx
func validatePhoneIndonesia(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return false
	}

	// Remove spaces and dashes
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	// Check for valid formats
	validPrefixes := []string{"08", "+628", "628"}
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(phone, prefix) {
			// Number after prefix should be 8-12 digits
			remaining := phone[len(prefix):]
			if len(remaining) >= 8 && len(remaining) <= 12 {
				return isNumeric(remaining)
			}
		}
	}
	return false
}

// validateStrongPassword validates password has uppercase, lowercase, number, and special char
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.Contains("!@#$%^&*()_+-=[]{}|;':\",./<>?", string(char)):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateDateRange validates end date is after start date
// Expects struct with StartDate and EndDate time.Time fields
func validateDateRange(fl validator.FieldLevel) bool {
	startField := fl.Field().FieldByName("StartDate")
	endField := fl.Field().FieldByName("EndDate")

	if !startField.IsValid() || !endField.IsValid() {
		return false
	}

	start := startField.Interface().(time.Time)
	end := endField.Interface().(time.Time)

	return end.After(start) || end.Equal(start)
}

// validateUniqueItems validates all items in a slice are unique
func validateUniqueItems(fl validator.FieldLevel) bool {
	slice := fl.Field()
	if slice.Kind() != reflect.Slice {
		return false
	}

	seen := make(map[interface{}]bool)
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i).Interface()
		if seen[item] {
			return false
		}
		seen[item] = true
	}
	return true
}
