package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type permissionUsecaseImpl struct {
	repo repositories.PermissionRepository
}

func NewPermissionUsecase(r repositories.PermissionRepository) PermissionUsecase {
	return &permissionUsecaseImpl{repo: r}
}

func (u *permissionUsecaseImpl) List(ctx context.Context) ([]*entities.Permission, error) {
	return u.repo.List(ctx)
}

func (u *permissionUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *permissionUsecaseImpl) Create(ctx context.Context, p *entities.Permission) error {
	if existing, err := u.repo.GetByResourceAction(ctx, p.Resource, p.Action); err == nil && existing != nil {
		return fmt.Errorf("permission already exists")
	} else if err != nil {
		return err
	}
	return u.repo.Create(ctx, p)
}

func (u *permissionUsecaseImpl) Update(ctx context.Context, p *entities.Permission) error {
	existing, err := u.repo.GetByID(ctx, p.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("permission not found")
	}
	if dup, err := u.repo.GetByResourceAction(ctx, p.Resource, p.Action); err == nil && dup != nil && dup.ID != p.ID {
		return fmt.Errorf("permission already exists")
	} else if err != nil {
		return err
	}
	return u.repo.Update(ctx, p)
}

func (u *permissionUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
