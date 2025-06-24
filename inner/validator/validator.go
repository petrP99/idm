package validator

import (
	"github.com/go-playground/validator/v10"
	"idm/inner/common"
	"strings"
)

// ValidationError represents structured validation errors
type ValidationError struct {
	FieldErrors map[string]string
}

func (e *ValidationError) Error() string {
	var errs []string
	for field, msg := range e.FieldErrors {
		errs = append(errs, field+": "+msg)
	}
	return "Validation failed: " + strings.Join(errs, ", ")
}

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate validates a request and returns structured error (RequestValidationError or other)
func (v *Validator) Validate(request any) error {
	err := v.validate.Struct(request)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Не валидационная ошибка — возвращаем как есть
		return err
	}

	// Преобразуем в кастомную ошибку RequestValidationError
	return common.MapValidationErrors(validationErrors)
}
