package usecases

import "context"

// AuthUsecase defines authentication business logic.
type AuthUsecase interface {
	Login(ctx context.Context, username, password string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Validate(ctx context.Context, token string) error
	Logout(ctx context.Context, token string) error
}
