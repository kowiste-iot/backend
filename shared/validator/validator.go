package validator

import (

	"backend/shared/errors"
	"reflect"
	"slices"
	"strings"
	"sync"
	roleDomain"backend/internal/features/user/domain"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	validate   *validator.Validate
	once       sync.Once
	validRoles []string
)

func InitValidator(roles []string) {
	validRoles = roles
	GetValidator()
}

// validateUUIDv7 validates if a string is a valid UUIDv7
func validateUUIDv7(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Handle empty string case if needed
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return false
	}

	return id.Version() == 7
}
func validateRoles(fl validator.FieldLevel) bool {
	roles, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}

	if len(roles) == 0 {
		return false
	}

	for _, rol := range roles {
		if !slices.ContainsFunc(roleDomain.AllRoles(validRoles), func(r roleDomain.Role) bool {
			return r.Name == rol
		}) {
			return false
		}
	}
	return true
}

func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		// Register function to get json tag name
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register UUIDv7 validator
		if err := validate.RegisterValidation("uuidv7", validateUUIDv7); err != nil {
			panic(err)
		}
		if err := validate.RegisterValidation("roles", validateRoles); err != nil {
			panic(err)
		}
	})
	return validate
}

// ValidateStruct validates a struct and returns AppError
func Validate(s interface{}) error {
	err := GetValidator().Struct(s)
	if err == nil {
		return nil
	}

	// Convert validator errors to AppError
	validationErrors := err.(validator.ValidationErrors)
	errorMessages := make([]string, 0, len(validationErrors))

	for _, e := range validationErrors {
		errorMessages = append(errorMessages, formatError(e))
	}

	return errors.NewValidation(
		strings.Join(errorMessages, "; "),
		err,
	)
}

// formatError formats a validation error into a readable message
func formatError(err validator.FieldError) string {
	field := err.Field()
	switch err.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email"
	case "min":
		return field + " must be at least " + err.Param()
	case "max":
		return field + " must be at most " + err.Param()
	case "uuidv7":
		return field + " must be a valid UUID version 7"
	default:
		return field + " failed on " + err.Tag() + " validation"
	}
}
