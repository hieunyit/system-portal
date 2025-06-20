package impl

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type inMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*entities.User
}

func NewUserRepository() repositories.UserRepository {
	return &inMemoryUserRepo{users: make(map[uuid.UUID]*entities.User)}
}

func (r *inMemoryUserRepo) Create(ctx context.Context, u *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *inMemoryUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.users[id], nil
}

func (r *inMemoryUserRepo) List(ctx context.Context) ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*entities.User, 0, len(r.users))
	for _, u := range r.users {
		out = append(out, u)
	}
	return out, nil
}

func (r *inMemoryUserRepo) Update(ctx context.Context, u *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

func (r *inMemoryUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, id)
	return nil
}
