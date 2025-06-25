package common

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

// RequestValidationError — ошибка валидации полей запроса
type RequestValidationError struct {
	FieldErrors map[string]string
}

// Обязательный метод для реализации интерфейса error
func (e RequestValidationError) Error() string {
	var errs []string
	for field, msg := range e.FieldErrors {
		errs = append(errs, fmt.Sprintf("%s: %s", field, msg))
	}
	return "Validation failed: " + strings.Join(errs, ", ")
}

// AlreadyExistsError — ресурс уже существует
type AlreadyExistsError struct {
	Resource string
	ID       any
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s with ID '%v' already exists", e.Resource, e.ID)
}

// NotFoundError — ресурс не найден
type NotFoundError struct {
	Resource string
	ID       any
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%v' not found", e.Resource, e.ID)
}

// InternalServerError — внутренняя ошибка сервера
type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return "Internal server error: " + e.Message
}

// MapValidationErrors — преобразует validator.ValidationErrors в RequestValidationError
func MapValidationErrors(validErrs validator.ValidationErrors) RequestValidationError {
	fieldErrors := make(map[string]string)

	for _, err := range validErrs {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		switch tag {
		case "required":
			fieldErrors[field] = "is required"
		case "min":
			fieldErrors[field] = fmt.Sprintf("must be at least %s characters", param)
		case "max":
			fieldErrors[field] = fmt.Sprintf("must not exceed %s characters", param)
		default:
			fieldErrors[field] = fmt.Sprintf("is invalid (%s)", tag)
		}
	}

	return RequestValidationError{FieldErrors: fieldErrors}
}
