package entities

import (
	"github.com/google/uuid"
	"time"
)

// SMTPConfig holds SMTP server configuration stored in the database
type SMTPConfig struct {
	ID        uuid.UUID
	Host      string
	Port      int
	Username  string
	Password  string
	From      string
	TLS       bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
