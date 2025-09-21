package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateUUID tests the ValidateUUID function
func TestValidateUUID(t *testing.T) {
	t.Run("Valid UUIDs should pass", func(t *testing.T) {
		validUUIDs := []string{
			"123e4567-e89b-12d3-a456-426614174000",
			"00000000-0000-0000-0000-000000000000",
			"ffffffff-ffff-ffff-ffff-ffffffffffff",
			"12345678-1234-1234-1234-123456789abc",
			"ABCDEFAB-ABCD-ABCD-ABCD-ABCDEFABCDEF", // uppercase hex
		}

		for _, uuid := range validUUIDs {
			err := ValidateUUID(uuid, "test_field")
			assert.NoError(t, err, "UUID %s should be valid", uuid)
		}
	})

	t.Run("Empty UUID should fail", func(t *testing.T) {
		err := ValidateUUID("", "user_id")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id is required")
	})

	t.Run("Invalid UUIDs should fail", func(t *testing.T) {
		invalidUUIDs := []string{
			"not-a-uuid",
			"123e4567-e89b-12d3-a456", // too short
			"123e4567-e89b-12d3-a456-426614174000-extra", // too long
			"123e4567-e89b-12d3-a456-42661417400g",       // invalid character 'g'
			"123e4567e89b12d3a456426614174000",           // missing dashes
			"123e4567-e89b-12d3-a456-42661417400",        // one char short
		}

		for _, uuid := range invalidUUIDs {
			err := ValidateUUID(uuid, "tool_id")

			assert.Error(t, err, "UUID %s should be invalid", uuid)
			assert.Contains(t, err.Error(), "tool_id must be a valid UUID")
		}
	})

	t.Run("Custom field names in error messages", func(t *testing.T) {
		err := ValidateUUID("", "custom_field")
		assert.Contains(t, err.Error(), "custom_field is required")

		err = ValidateUUID("invalid", "another_field")
		assert.Contains(t, err.Error(), "another_field must be a valid UUID")
	})
}
