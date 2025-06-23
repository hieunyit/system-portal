package usecases

import (
	"context"
	"fmt"
	"time"

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
	existing, err := g.repo.GetByName(ctx, gr.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("group already exists")
	}
	return g.repo.Create(ctx, gr)
}

func (g *groupUsecaseImpl) List(ctx context.Context, filter *entities.GroupFilter) ([]*entities.PortalGroup, int, error) {
	groups, total, err := g.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	if g.permRepo != nil {
		for _, grp := range groups {
			perms, _ := g.permRepo.GetByGroup(ctx, grp.ID)
			grp.Permissions = perms
		}
	}
	return groups, total, nil
}

func (g *groupUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (*entities.PortalGroup, error) {
	grp, err := g.repo.GetByID(ctx, id)
	if err != nil || grp == nil {
		return grp, err
	}
	if g.permRepo != nil {
		grp.Permissions, _ = g.permRepo.GetByGroup(ctx, grp.ID)
	}
	return grp, nil
}

func (g *groupUsecaseImpl) GetByName(ctx context.Context, name string) (*entities.PortalGroup, error) {
	return g.repo.GetByName(ctx, name)
}

func (g *groupUsecaseImpl) Update(ctx context.Context, gr *entities.PortalGroup) error {
	existing, err := g.repo.GetByID(ctx, gr.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("group not found")
	}
	if other, err := g.repo.GetByName(ctx, gr.Name); err == nil && other != nil && other.ID != gr.ID {
		return fmt.Errorf("group already exists")
	}
	gr.UpdatedAt = time.Now()
	return g.repo.Update(ctx, gr)
}

func (g *groupUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := g.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("group not found")
	}
	return g.repo.Delete(ctx, id)
}

func (g *groupUsecaseImpl) UpdatePermissions(ctx context.Context, id uuid.UUID, permIDs []uuid.UUID) error {
	if g.permRepo == nil {
		return nil
	}
	grp, err := g.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if grp == nil {
		return fmt.Errorf("group not found")
	}
	for _, pid := range permIDs {
		p, err := g.permRepo.GetByID(ctx, pid)
		if err != nil {
			return err
		}
		if p == nil {
			return fmt.Errorf("permission not found")
		}
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
