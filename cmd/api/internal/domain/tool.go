package domain

import (
	"fmt"
	"time"
)

type ToolStatus string

const (
	ToolStatusInOffice    ToolStatus = "IN_OFFICE"
	ToolStatusCheckedOut  ToolStatus = "CHECKED_OUT"
	ToolStatusMaintenance ToolStatus = "MAINTENANCE"
	ToolStatusLost        ToolStatus = "LOST"
)

func (s ToolStatus) IsValid() bool {
	switch s {
	case ToolStatusInOffice, ToolStatusCheckedOut, ToolStatusLost, ToolStatusMaintenance:
		return true
	default:
		return false
	}
}

type Tool struct {
	ID               *string    `json:"id"`
	Name             string     `json:"name"`
	Status           ToolStatus `json:"status"`
	CurrentUserId    *string    `json:"current_user_id,omitempty"`
	LastCheckedOutAt *time.Time `json:"last_checked_out_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func NewTool(name string, status ToolStatus) (Tool, error) {
	if status == "" {
		status = ToolStatusInOffice
	}
	t := Tool{Name: name, Status: status}

	return t, t.Validate()
}

func ValidateToolStatus(status ToolStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("%w: invalid status %s", ErrValidation, status)
	}

	return nil
}
func (t *Tool) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if err := ValidateToolStatus(t.Status); err != nil {
		return err
	}

	return nil
}
