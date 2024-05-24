package middlewares

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func FormatValidationError(err error) map[string]string {
	errorsMap := make(map[string]string)
	if _, ok := err.(*validator.InvalidValidationError); ok {
		errorsMap["error"] = err.Error()
		return errorsMap
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			errorsMap[fieldErr.Field()] = "failed on the '" + fieldErr.Tag() + "' tag"
		}
	}

	return errorsMap
}
