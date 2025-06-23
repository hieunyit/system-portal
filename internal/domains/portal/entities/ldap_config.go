package entities

import (
	"github.com/google/uuid"
	"time"
)

// LDAPConfig stores connection info for an LDAP server
type LDAPConfig struct {
	ID           uuid.UUID
	Host         string
	Port         int
	BindDN       string
	BindPassword string
	BaseDN       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
