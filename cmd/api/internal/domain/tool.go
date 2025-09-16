package domain

import "time"

type Tool struct {
	ID        string
	Name      string
	Status    string
	CreatedAt time.Time
}

func (t *Tool) Validate() error {
	return nil
}
