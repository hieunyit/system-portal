package usecases

import (
	"context"
	"system-portal/internal/domains/portal/entities"

	"github.com/google/uuid"
)

type UserUsecase interface {
	Create(ctx context.Context, u *entities.PortalUser) error
	List(ctx context.Context, filter *entities.UserFilter) ([]*entities.PortalUser, int, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.PortalUser, error)
	Update(ctx context.Context, u *entities.PortalUser) error
	Delete(ctx context.Context, id uuid.UUID) error
}
