package entities

import (
	"time"

	"github.com/google/uuid"
)

// User represents a portal user used for authentication.
type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	GroupID      uuid.UUID
	CreatedAt    time.Time
}
