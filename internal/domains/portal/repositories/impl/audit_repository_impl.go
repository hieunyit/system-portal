package impl

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type inMemoryAuditRepo struct {
	mu   sync.RWMutex
	logs map[uuid.UUID]*entities.AuditLog
}

func NewAuditRepository() repositories.AuditRepository {
	return &inMemoryAuditRepo{logs: make(map[uuid.UUID]*entities.AuditLog)}
}

func (r *inMemoryAuditRepo) Add(ctx context.Context, l *entities.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs[l.ID] = l
	return nil
}

func (r *inMemoryAuditRepo) List(ctx context.Context) ([]*entities.AuditLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*entities.AuditLog, 0, len(r.logs))
	for _, l := range r.logs {
		out = append(out, l)
	}
	return out, nil
}

func (r *inMemoryAuditRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.logs[id], nil
}
