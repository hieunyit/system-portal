package entities

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog records portal activities.
type AuditLog struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Action    string
	Resource  string
	Success   bool
	CreatedAt time.Time
}
