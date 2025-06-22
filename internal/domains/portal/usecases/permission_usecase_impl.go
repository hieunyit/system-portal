package usecases

import (
	"context"
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
	return u.repo.Create(ctx, p)
}

func (u *permissionUsecaseImpl) Update(ctx context.Context, p *entities.Permission) error {
	return u.repo.Update(ctx, p)
}

func (u *permissionUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
