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

func (t *Tool) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}
	switch t.Status {
	case ToolStatusInOffice, ToolStatusCheckedOut, ToolStatusLost, ToolStatusMaintenance:
		//Valid
	default:
		return fmt.Errorf("invalid status %s", t.Status)
	}

	return nil
}
