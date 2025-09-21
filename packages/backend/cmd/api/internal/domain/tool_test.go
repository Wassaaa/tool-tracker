package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestToolStatus_IsValid tests the ToolStatus validation
func TestToolStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   ToolStatus
		expected bool
	}{
		{"Valid IN_OFFICE", ToolStatusInOffice, true},
		{"Valid CHECKED_OUT", ToolStatusCheckedOut, true},
		{"Valid MAINTENANCE", ToolStatusMaintenance, true},
		{"Valid LOST", ToolStatusLost, true},
		{"Invalid empty", ToolStatus(""), false},
		{"Invalid random", ToolStatus("RANDOM_STATUS"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNewTool tests tool creation
func TestNewTool(t *testing.T) {
	t.Run("Valid tool creation", func(t *testing.T) {
		tool, err := NewTool("Hammer", ToolStatusInOffice)

		require.NoError(t, err)
		assert.Equal(t, "Hammer", tool.Name)
		assert.Equal(t, ToolStatusInOffice, tool.Status)
	})

	t.Run("Empty name should fail", func(t *testing.T) {
		_, err := NewTool("", ToolStatusInOffice)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("Invalid status should fail", func(t *testing.T) {
		_, err := NewTool("Hammer", ToolStatus("INVALID"))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})

	t.Run("Default status when empty", func(t *testing.T) {
		tool, err := NewTool("Hammer", "")

		require.NoError(t, err)
		assert.Equal(t, ToolStatusInOffice, tool.Status)
	})
}

// TestTool_Validate tests the Tool Validate method comprehensively
func TestTool_Validate(t *testing.T) {
	t.Run("Valid tool should pass validation", func(t *testing.T) {
		tool := Tool{
			Name:   "Hammer",
			Status: ToolStatusInOffice,
		}

		err := tool.Validate()
		assert.NoError(t, err)
	})

	t.Run("Empty name should fail validation", func(t *testing.T) {
		tool := Tool{
			Name:   "",
			Status: ToolStatusInOffice,
		}

		err := tool.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("All valid statuses should pass validation", func(t *testing.T) {
		validStatuses := []ToolStatus{
			ToolStatusInOffice,
			ToolStatusCheckedOut,
			ToolStatusMaintenance,
			ToolStatusLost,
		}

		for _, status := range validStatuses {
			tool := Tool{
				Name:   "Test Tool",
				Status: status,
			}

			err := tool.Validate()
			assert.NoError(t, err, "Status %s should be valid", status)
		}
	})

	t.Run("Invalid status should fail validation", func(t *testing.T) {
		invalidStatuses := []ToolStatus{
			ToolStatus(""),
			ToolStatus("INVALID"),
			ToolStatus("random_status"),
			ToolStatus("in_office"), // case sensitive
		}

		for _, status := range invalidStatuses {
			tool := Tool{
				Name:   "Test Tool",
				Status: status,
			}

			err := tool.Validate()
			assert.Error(t, err, "Status %s should be invalid", status)
			assert.Contains(t, err.Error(), "invalid status")
		}
	})

	t.Run("Multiple validation errors - should fail on first", func(t *testing.T) {
		tool := Tool{
			Name:   "",                    // First error
			Status: ToolStatus("INVALID"), // Second error
		}

		err := tool.Validate()
		assert.Error(t, err)
		// Should fail on first validation error (name)
		assert.Contains(t, err.Error(), "name is required")
	})
}
