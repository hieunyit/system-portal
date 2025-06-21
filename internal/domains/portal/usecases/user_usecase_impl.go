package usecases

import (
	"context"

	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecaseImpl struct{ repo repositories.UserRepository }

func NewUserUsecase(repo repositories.UserRepository) UserUsecase {
	return &userUsecaseImpl{repo: repo}
}

func (u *userUsecaseImpl) Create(ctx context.Context, user *entities.PortalUser) error {
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}
	return u.repo.Create(ctx, user)
}

func (u *userUsecaseImpl) List(ctx context.Context) ([]*entities.PortalUser, error) {
	return u.repo.List(ctx)
}

func (u *userUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (*entities.PortalUser, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *userUsecaseImpl) Update(ctx context.Context, user *entities.PortalUser) error {
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}
	return u.repo.Update(ctx, user)
}

func (u *userUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
