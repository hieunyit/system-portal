package impl

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"system-portal/internal/domains/auth/entities"
	"system-portal/internal/domains/auth/repositories"
)

// InMemorySessionRepository is a naive in-memory session store used for demos.
type InMemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[uuid.UUID]*entities.Session
}

// NewSessionRepository creates a new repository instance.
func NewSessionRepository() repositories.SessionRepository {
	return &InMemorySessionRepository{sessions: make(map[uuid.UUID]*entities.Session)}
}

func (r *InMemorySessionRepository) Create(ctx context.Context, s *entities.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s.ID] = s
	return nil
}

func (r *InMemorySessionRepository) GetByTokenHash(ctx context.Context, hash string) (*entities.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.sessions {
		if s.TokenHash == hash && s.IsActive {
			return s, nil
		}
	}
	return nil, nil
}

func (r *InMemorySessionRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.sessions[id]; ok {
		s.IsActive = false
	}
	return nil
}
