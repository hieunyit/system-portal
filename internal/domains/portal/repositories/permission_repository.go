package repositories

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type PermissionRepository interface {
	List(ctx context.Context) ([]*entities.Permission, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error)
	Create(ctx context.Context, p *entities.Permission) error
	Update(ctx context.Context, p *entities.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByGroup(ctx context.Context, groupID uuid.UUID) ([]*entities.Permission, error)
	SetForGroup(ctx context.Context, groupID uuid.UUID, permIDs []uuid.UUID) error
	HasGroupPermission(ctx context.Context, groupName, resource, action string) (bool, error)
}
