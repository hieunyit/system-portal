package repositories

import (
	"context"
	"system-portal/internal/domains/portal/entities"
)

type EmailTemplateRepository interface {
	GetByAction(ctx context.Context, action string) (*entities.EmailTemplate, error)
	Upsert(ctx context.Context, tpl *entities.EmailTemplate) error
}
