package usecases

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type userUsecaseImpl struct{ repo repositories.UserRepository }

func NewUserUsecase(repo repositories.UserRepository) UserUsecase {
	return &userUsecaseImpl{repo: repo}
}

func (u *userUsecaseImpl) Create(ctx context.Context, user *entities.User) error {
	return u.repo.Create(ctx, user)
}

func (u *userUsecaseImpl) List(ctx context.Context) ([]*entities.User, error) {
	return u.repo.List(ctx)
}

func (u *userUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUsecaseImpl) Update(ctx context.Context, user *entities.User) error {
	return u.repo.Update(ctx, user)
}

func (u *userUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
