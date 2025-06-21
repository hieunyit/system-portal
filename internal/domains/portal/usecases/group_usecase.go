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
}
