package usecases

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type groupUsecaseImpl struct{ repo repositories.GroupRepository }

func NewGroupUsecase(repo repositories.GroupRepository) GroupUsecase {
	return &groupUsecaseImpl{repo: repo}
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
