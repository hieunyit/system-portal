package repositories

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entities.PortalGroup) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.PortalGroup, error)
	GetByName(ctx context.Context, name string) (*entities.PortalGroup, error)
       List(ctx context.Context, filter *entities.GroupFilter) ([]*entities.PortalGroup, int, error)
	Update(ctx context.Context, group *entities.PortalGroup) error
	Delete(ctx context.Context, id uuid.UUID) error
}
