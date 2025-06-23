package entities

import (
	"github.com/google/uuid"
	"time"
)

// OpenVPNConfig stores connection info for the OpenVPN XML-RPC service
type OpenVPNConfig struct {
	ID        uuid.UUID
	Host      string
	Username  string
	Password  string
	Port      int
	CreatedAt time.Time
	UpdatedAt time.Time
}
