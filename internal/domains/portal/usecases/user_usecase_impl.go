package usecases

import (
	"context"
	"fmt"

	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecaseImpl struct {
	repo      repositories.UserRepository
	groupRepo repositories.GroupRepository
}

func NewUserUsecase(repo repositories.UserRepository, group repositories.GroupRepository) UserUsecase {
	return &userUsecaseImpl{repo: repo, groupRepo: group}
}

func (u *userUsecaseImpl) Create(ctx context.Context, user *entities.PortalUser) error {
	if existing, _ := u.repo.GetByUsername(ctx, user.Username); existing != nil {
		return fmt.Errorf("username already exists")
	}
	if existing, _ := u.repo.GetByEmail(ctx, user.Email); existing != nil {
		return fmt.Errorf("email already exists")
	}
	if g, _ := u.groupRepo.GetByID(ctx, user.GroupID); g == nil {
		return fmt.Errorf("group not found")
	}
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
	existing, err := u.repo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("user not found")
	}
	if other, _ := u.repo.GetByUsername(ctx, user.Username); other != nil && other.ID != user.ID {
		return fmt.Errorf("username already exists")
	}
	if other, _ := u.repo.GetByEmail(ctx, user.Email); other != nil && other.ID != user.ID {
		return fmt.Errorf("email already exists")
	}
	if g, _ := u.groupRepo.GetByID(ctx, user.GroupID); g == nil {
		return fmt.Errorf("group not found")
	}
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
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("user not found")
	}
	return u.repo.Delete(ctx, id)
}
