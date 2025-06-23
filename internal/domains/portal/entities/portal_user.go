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

// UserFilter defines optional filters and pagination for listing users.
type UserFilter struct {
	Username string
	Email    string
	GroupID  uuid.UUID
	Page     int
	Limit    int
	Offset   int
}

// SetDefaults ensures pagination defaults and calculates the offset.
func (f *UserFilter) SetDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	f.Offset = (f.Page - 1) * f.Limit
}
