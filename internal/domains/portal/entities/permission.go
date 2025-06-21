package entities

import "github.com/google/uuid"

// Permission represents a single allowed action.
type Permission struct {
	ID          uuid.UUID
	Resource    string
	Action      string
	Description string
}
