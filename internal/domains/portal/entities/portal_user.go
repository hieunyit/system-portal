package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a portal user entity.
type PortalUser struct {
	ID        uuid.UUID
	Username  string
	Email     string
	FullName  string
	Password  string
	GroupID   uuid.UUID
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
