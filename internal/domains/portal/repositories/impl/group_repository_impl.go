package impl

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type inMemoryGroupRepo struct {
	mu     sync.RWMutex
	groups map[uuid.UUID]*entities.Group
}

func NewGroupRepository() repositories.GroupRepository {
	return &inMemoryGroupRepo{groups: make(map[uuid.UUID]*entities.Group)}
}

func (r *inMemoryGroupRepo) Create(ctx context.Context, g *entities.Group) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups[g.ID] = g
	return nil
}

func (r *inMemoryGroupRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.groups[id], nil
}

func (r *inMemoryGroupRepo) List(ctx context.Context) ([]*entities.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*entities.Group, 0, len(r.groups))
	for _, g := range r.groups {
		out = append(out, g)
	}
	return out, nil
}
