package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"system-portal/internal/domains/auth/entities"
	"system-portal/internal/domains/auth/repositories"
	repoimpl "system-portal/internal/domains/auth/repositories/impl"
	"system-portal/pkg/jwt"
)

type authUsecaseImpl struct {
	sessions repositories.SessionRepository
	jwt      *jwt.RSAService
}

func NewAuthUsecase(jwtSvc *jwt.RSAService) AuthUsecase {
	return &authUsecaseImpl{sessions: repoimpl.NewSessionRepository(), jwt: jwtSvc}
}

func (u *authUsecaseImpl) Login(ctx context.Context, username, password string) (string, string, error) {
	access, _ := u.jwt.GenerateAccessToken(username, "admin")
	refresh, _ := u.jwt.GenerateRefreshToken(username, "admin")
	s := &entities.Session{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		TokenHash:        access,
		RefreshTokenHash: refresh,
		ExpiresAt:        time.Now().Add(time.Hour),
		RefreshExpiresAt: time.Now().Add(24 * time.Hour),
		IsActive:         true,
		CreatedAt:        time.Now(),
	}
	u.sessions.Create(ctx, s)
	return access, refresh, nil
}

func (u *authUsecaseImpl) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := u.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	return u.Login(ctx, claims.Username, "")
}

func (u *authUsecaseImpl) Validate(ctx context.Context, token string) error {
	_, err := u.jwt.ValidateAccessToken(token)
	return err
}

func (u *authUsecaseImpl) Logout(ctx context.Context, token string) error {
	// In-memory implementation simply ignores logout
	return nil
}
