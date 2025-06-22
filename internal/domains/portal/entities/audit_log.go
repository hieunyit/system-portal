package entities

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog records portal activities.
type AuditLog struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Username     string
	UserGroup    string
	Action       string
	ResourceType string
	ResourceName string
	IPAddress    string
	Success      bool
	CreatedAt    time.Time
}
