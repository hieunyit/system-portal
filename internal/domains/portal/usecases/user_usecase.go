package usecases

import (
	"context"
	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type UserUsecase interface {
	Create(ctx context.Context, u *entities.User) error
	List(ctx context.Context) ([]*entities.User, error)
	Get(ctx context.Context, id uuid.UUID) (*entities.User, error)
	Update(ctx context.Context, u *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
