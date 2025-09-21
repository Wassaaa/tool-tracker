package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func respondDomainError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	body := apiError{Code: "internal_error", Message: "internal error"}

	switch {
	case errors.Is(err, domain.ErrValidation):
		status = http.StatusBadRequest
		body = apiError{Code: "validation_error", Message: err.Error()}
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
		body = apiError{Code: "conflict", Message: err.Error()}
	case errors.Is(err, domain.ErrToolNotFound):
		status = http.StatusNotFound
		body = apiError{Code: "tool_not_found", Message: err.Error()}
	case errors.Is(err, domain.ErrUserNotFound):
		status = http.StatusNotFound
		body = apiError{Code: "user_not_found", Message: err.Error()}
	case errors.Is(err, domain.ErrEventNotFound):
		status = http.StatusNotFound
		body = apiError{Code: "event_not_found", Message: err.Error()}
	}

	c.JSON(status, gin.H{"error": body})
}

// validationErr wraps domain.ErrValidation with a contextual message (field + detail).
func validationErr(field, msg string) error {
	cleanField := strings.TrimSpace(field)
	if cleanField == "" {
		return fmt.Errorf("%w: %s", domain.ErrValidation, msg)
	}
	return fmt.Errorf("%w: %s %s", domain.ErrValidation, cleanField, msg)
}
