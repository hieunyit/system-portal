package repositories

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entities.PortalGroup) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.PortalGroup, error)
	List(ctx context.Context) ([]*entities.PortalGroup, error)
}
