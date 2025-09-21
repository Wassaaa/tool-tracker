package domain

import (
	"fmt"
	"regexp"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ValidateUUID validates a UUID string and returns an ErrValidation-wrapped error if invalid.
func ValidateUUID(id string, field string) error {
	if id == "" {
		return fmt.Errorf("%w: %s is required", ErrValidation, field)
	}
	if !uuidRegex.MatchString(id) {
		return fmt.Errorf("%w: %s must be a valid UUID", ErrValidation, field)
	}
	return nil
}
