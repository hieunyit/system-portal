package usecases

import (
	"context"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type AuditUsecase interface {
	Add(ctx context.Context, log *entities.AuditLog) error
	List(ctx context.Context, filter *entities.AuditFilter) ([]*entities.AuditLog, int, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error)
}
