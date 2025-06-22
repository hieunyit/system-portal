package usecases

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type groupUsecaseImpl struct {
	repo     repositories.GroupRepository
	permRepo repositories.PermissionRepository
}

func NewGroupUsecase(repo repositories.GroupRepository, perm repositories.PermissionRepository) GroupUsecase {
	return &groupUsecaseImpl{repo: repo, permRepo: perm}
}

func (g *groupUsecaseImpl) Create(ctx context.Context, gr *entities.PortalGroup) error {
	return g.repo.Create(ctx, gr)
}

func (g *groupUsecaseImpl) List(ctx context.Context) ([]*entities.PortalGroup, error) {
	return g.repo.List(ctx)
}

func (g *groupUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (*entities.PortalGroup, error) {
	return g.repo.GetByID(ctx, id)
}

func (g *groupUsecaseImpl) GetByName(ctx context.Context, name string) (*entities.PortalGroup, error) {
	return g.repo.GetByName(ctx, name)
}

func (g *groupUsecaseImpl) UpdatePermissions(ctx context.Context, id uuid.UUID, permIDs []uuid.UUID) error {
	if g.permRepo == nil {
		return nil
	}
	return g.permRepo.SetForGroup(ctx, id, permIDs)
}

func (g *groupUsecaseImpl) GetPermissions(ctx context.Context, id uuid.UUID) ([]*entities.Permission, error) {
	if g.permRepo == nil {
		return nil, nil
	}
	return g.permRepo.GetByGroup(ctx, id)
}

func (g *groupUsecaseImpl) ListPermissions(ctx context.Context) ([]*entities.Permission, error) {
	if g.permRepo == nil {
		return nil, nil
	}
	return g.permRepo.List(ctx)
}
