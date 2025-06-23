package repositories

import (
	"context"

	"system-portal/internal/domains/portal/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.PortalUser) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.PortalUser, error)
	GetByUsername(ctx context.Context, username string) (*entities.PortalUser, error)
	GetByEmail(ctx context.Context, email string) (*entities.PortalUser, error)
       List(ctx context.Context, filter *entities.UserFilter) ([]*entities.PortalUser, int, error)
	Update(ctx context.Context, user *entities.PortalUser) error
	Delete(ctx context.Context, id uuid.UUID) error
}
