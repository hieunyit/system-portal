package usecases

import (
	"context"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type PermissionUsecase interface {
	List(ctx context.Context) ([]*entities.Permission, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Permission, error)
	Create(ctx context.Context, p *entities.Permission) error
	Update(ctx context.Context, p *entities.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
}
