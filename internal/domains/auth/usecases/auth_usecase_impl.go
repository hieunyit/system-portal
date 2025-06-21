package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"system-portal/internal/domains/auth/entities"
	"system-portal/internal/domains/auth/repositories"
	portalrepos "system-portal/internal/domains/portal/repositories"
	"system-portal/pkg/jwt"
	"system-portal/pkg/logger"
)

type authUsecaseImpl struct {
	sessions repositories.SessionRepository
	users    portalrepos.UserRepository
	jwt      *jwt.RSAService
}

func NewAuthUsecase(sessionRepo repositories.SessionRepository, userRepo portalrepos.UserRepository, jwtSvc *jwt.RSAService) AuthUsecase {
	return &authUsecaseImpl{sessions: sessionRepo, users: userRepo, jwt: jwtSvc}
}

func (u *authUsecaseImpl) Login(ctx context.Context, username, password string) (string, string, error) {
	usr, err := u.users.GetByUsername(ctx, username)
	if err != nil {
		logger.Log.WithError(err).Error("failed to fetch user")
		return "", "", errors.New("invalid credentials")
	}
	if usr == nil || !usr.IsActive {
		logger.Log.WithField("username", username).Warn("user not found or inactive")
		return "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		logger.Log.WithField("username", username).Warn("password mismatch")
		return "", "", errors.New("invalid credentials")
	}

	role := "support"
	if username == "admin" {
		role = "admin"
	}
	access, _ := u.jwt.GenerateAccessToken(username, role)
	refresh, _ := u.jwt.GenerateRefreshToken(username, role)
	s := &entities.Session{
		ID:               uuid.New(),
		UserID:           usr.ID,
		TokenHash:        access,
		RefreshTokenHash: refresh,
		ExpiresAt:        time.Now().Add(time.Hour),
		RefreshExpiresAt: time.Now().Add(24 * time.Hour),
		IsActive:         true,
		CreatedAt:        time.Now(),
	}
	u.sessions.Create(ctx, s)
	logger.Log.WithField("username", username).Info("session created")
	return access, refresh, nil
}

func (u *authUsecaseImpl) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := u.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.Log.WithError(err).Warn("refresh token validation failed")
		return "", "", err
	}
	logger.Log.WithField("username", claims.Username).Info("refresh token validated")
	return u.Login(ctx, claims.Username, "")
}

func (u *authUsecaseImpl) Validate(ctx context.Context, token string) error {
	_, err := u.jwt.ValidateAccessToken(token)
	if err != nil {
		logger.Log.WithError(err).Warn("access token invalid")
	}
	return err
}

func (u *authUsecaseImpl) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	if _, err := u.jwt.ValidateAccessToken(token); err != nil {
		logger.Log.WithError(err).Warn("logout token validation failed")
		return err
	}
	sess, err := u.sessions.GetByTokenHash(ctx, token)
	if err != nil {
		logger.Log.WithError(err).Error("failed to get session by token")
		return err
	}
	if sess == nil {
		return nil
	}
	logger.Log.WithField("sessionID", sess.ID).Info("deactivating session")
	return u.sessions.Deactivate(ctx, sess.ID)
}
