package entities

import (
	"time"

	"github.com/google/uuid"
)

// Group represents a portal group with permissions.
type PortalGroup struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Permissions []*Permission
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GroupFilter defines optional filters and pagination for listing groups.
type GroupFilter struct {
	Name   string
	Page   int
	Limit  int
	Offset int
}

// SetDefaults ensures pagination defaults and calculates the offset.
func (f *GroupFilter) SetDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	f.Offset = (f.Page - 1) * f.Limit
}

// Backward compatibility alias
type Group = PortalGroup
