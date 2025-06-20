package usecases

import (
	"context"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type GroupUsecase interface {
	Create(ctx context.Context, g *entities.Group) error
	List(ctx context.Context) ([]*entities.Group, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.Group, error)
}
