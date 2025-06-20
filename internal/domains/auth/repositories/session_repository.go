package repositories

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/auth/entities"
)

// SessionRepository stores active sessions.
type SessionRepository interface {
	Create(ctx context.Context, s *entities.Session) error
	GetByTokenHash(ctx context.Context, hash string) (*entities.Session, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}
