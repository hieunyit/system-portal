package usecases

import (
	"context"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type GroupUsecase interface {
	Create(ctx context.Context, g *entities.PortalGroup) error
	List(ctx context.Context) ([]*entities.PortalGroup, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.PortalGroup, error)
	GetByName(ctx context.Context, name string) (*entities.PortalGroup, error)
	Update(ctx context.Context, g *entities.PortalGroup) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePermissions(ctx context.Context, id uuid.UUID, permIDs []uuid.UUID) error
	GetPermissions(ctx context.Context, id uuid.UUID) ([]*entities.Permission, error)
	ListPermissions(ctx context.Context) ([]*entities.Permission, error)
}
