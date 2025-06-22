package repositories

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
)

type PermissionRepository interface {
	List(ctx context.Context) ([]*entities.Permission, error)
	GetByGroup(ctx context.Context, groupID uuid.UUID) ([]*entities.Permission, error)
	SetForGroup(ctx context.Context, groupID uuid.UUID, permIDs []uuid.UUID) error
	HasGroupPermission(ctx context.Context, groupName, resource, action string) (bool, error)
}
