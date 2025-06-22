package repositories

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type AuditRepository interface {
	Add(ctx context.Context, log *entities.AuditLog) error
	List(ctx context.Context, filter *entities.AuditFilter) ([]*entities.AuditLog, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error)
}
