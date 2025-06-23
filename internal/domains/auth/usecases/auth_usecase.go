package usecases

import (
	"context"

	"github.com/google/uuid"
)

// AuthUsecase defines authentication business logic.
type AuthUsecase interface {
	Login(ctx context.Context, username, password, ip string) (string, string, uuid.UUID, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Validate(ctx context.Context, token string) error
	Logout(ctx context.Context, token string) error
}
