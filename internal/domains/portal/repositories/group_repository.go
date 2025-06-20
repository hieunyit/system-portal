package repositories

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type GroupRepository interface {
	Create(ctx context.Context, group *entities.Group) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Group, error)
	List(ctx context.Context) ([]*entities.Group, error)
}
