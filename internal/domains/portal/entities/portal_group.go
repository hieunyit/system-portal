package entities

import (
	"time"

	"github.com/google/uuid"
)

// Group represents a portal group with permissions.
type Group struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Permissions []Permission
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
