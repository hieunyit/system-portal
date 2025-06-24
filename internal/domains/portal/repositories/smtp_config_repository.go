package repositories

import (
	"context"
	"system-portal/internal/domains/portal/entities"
)

type SMTPConfigRepository interface {
	Get(ctx context.Context) (*entities.SMTPConfig, error)
	Create(ctx context.Context, cfg *entities.SMTPConfig) error
	Update(ctx context.Context, cfg *entities.SMTPConfig) error
	Delete(ctx context.Context) error
}
