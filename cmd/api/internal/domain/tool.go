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

type Tool struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Status    ToolStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewTool(name string, status ToolStatus) (Tool, error) {
    if status == "" {
        status = ToolStatusInOffice
    }
    t := Tool{Name: name, Status: status}
    return t, t.Validate()
}

func ValidateToolStatus(status ToolStatus) error {
	switch status {
	case ToolStatusInOffice, ToolStatusCheckedOut, ToolStatusLost, ToolStatusMaintenance:
		return nil
	default:
		return fmt.Errorf("invalid status %s", status)
	}
}
func (t *Tool) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}
	if err := ValidateToolStatus(t.Status); err != nil {
		return err
	}
	return nil
}
