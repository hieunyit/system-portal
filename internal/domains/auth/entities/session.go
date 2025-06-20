package entities

import (
	"time"

	"github.com/google/uuid"
)

// Session represents a user login session.
type Session struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	TokenHash        string
	RefreshTokenHash string
	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
	IsActive         bool
	CreatedAt        time.Time
}
