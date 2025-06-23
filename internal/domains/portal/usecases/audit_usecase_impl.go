package usecases

import (
	"context"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type auditUsecaseImpl struct{ repo repositories.AuditRepository }

func NewAuditUsecase(repo repositories.AuditRepository) AuditUsecase {
	return &auditUsecaseImpl{repo: repo}
}

func (a *auditUsecaseImpl) Add(ctx context.Context, l *entities.AuditLog) error {
	return a.repo.Add(ctx, l)
}

func (a *auditUsecaseImpl) List(ctx context.Context, f *entities.AuditFilter) ([]*entities.AuditLog, int, error) {
	return a.repo.List(ctx, f)
}

func (a *auditUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error) {
	return a.repo.GetByID(ctx, id)
}
