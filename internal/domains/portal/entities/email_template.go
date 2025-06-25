package entities

import (
	"github.com/google/uuid"
	"time"
)

// EmailTemplate stores customizable email templates for various actions
// Action is a unique identifier like "create_user" or "expiration"
type EmailTemplate struct {
	ID        uuid.UUID
	Action    string
	Subject   string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
