package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

// respondDomainError maps domain-level sentinel errors to HTTP responses.
// Falls back to 500 for unknown errors.
func respondDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrToolNotFound), errors.Is(err, domain.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
