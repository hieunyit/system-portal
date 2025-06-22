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
	"system-portal/pkg/utils"
)

type authUsecaseImpl struct {
	sessions repositories.SessionRepository
	users    portalrepos.UserRepository
	jwt      *jwt.RSAService
}

func NewAuthUsecase(sessionRepo repositories.SessionRepository, userRepo portalrepos.UserRepository, jwtSvc *jwt.RSAService) AuthUsecase {
	return &authUsecaseImpl{sessions: sessionRepo, users: userRepo, jwt: jwtSvc}
}

func (u *authUsecaseImpl) Login(ctx context.Context, username, password, ip string) (string, string, uuid.UUID, string, error) {
	logger.Log.WithField("username", username).Info("login attempt")
	usr, err := u.users.GetByUsername(ctx, username)
	if err != nil {
		logger.Log.WithError(err).Error("failed to fetch user")
		return "", "", uuid.Nil, "", errors.New("invalid credentials")
	}
	if usr == nil {
		logger.Log.WithField("username", username).Warn("user not found")
		return "", "", uuid.Nil, "", errors.New("invalid credentials")
	}
	if !usr.IsActive {
		logger.Log.WithField("username", username).Warn("user inactive")
		return "", "", uuid.Nil, "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password)); err != nil {
		logger.Log.WithField("username", username).Warn("password mismatch")
		return "", "", uuid.Nil, "", errors.New("invalid credentials")
	}

	if cost, err := bcrypt.Cost([]byte(usr.Password)); err == nil && cost > bcrypt.DefaultCost {
		if newHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err == nil {
			usr.Password = string(newHash)
			if err := u.users.Update(ctx, usr); err != nil {
				logger.Log.WithError(err).Warn("failed to update password hash")
			} else {
				logger.Log.WithField("username", username).Debug("password rehashed with lower cost")
			}
		}
	}

	role := "support"
	if username == "admin" {
		role = "admin"
	}
	access, _ := u.jwt.GenerateAccessToken(username, role)
	refresh, _ := u.jwt.GenerateRefreshToken(username, role)
	accessHash := utils.HashString(access)
	refreshHash := utils.HashString(refresh)
	s := &entities.Session{
		ID:               uuid.New(),
		UserID:           usr.ID,
		TokenHash:        accessHash,
		RefreshTokenHash: refreshHash,
		ExpiresAt:        time.Now().Add(time.Hour),
		RefreshExpiresAt: time.Now().Add(24 * time.Hour),
		IsActive:         true,
		CreatedAt:        time.Now(),
		IPAddress:        ip,
	}
	u.sessions.Create(ctx, s)
	logger.Log.WithField("username", username).Info("session created")
	return access, refresh, usr.ID, role, nil
}

func (u *authUsecaseImpl) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := u.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.Log.WithError(err).Warn("refresh token validation failed")
		return "", "", err
	}
	logger.Log.WithField("username", claims.Username).Info("refresh token validated")

	usr, err := u.users.GetByUsername(ctx, claims.Username)
	if err != nil {
		logger.Log.WithError(err).Error("failed to fetch user for refresh")
		return "", "", err
	}
	if usr == nil || !usr.IsActive {
		logger.Log.WithField("username", claims.Username).Warn("user not found or inactive")
		return "", "", errors.New("invalid credentials")
	}

	role := "support"
	if claims.Username == "admin" {
		role = "admin"
	}
	access, _ := u.jwt.GenerateAccessToken(claims.Username, role)
	refreshNew, _ := u.jwt.GenerateRefreshToken(claims.Username, role)
	accessHash := utils.HashString(access)
	refreshHash := utils.HashString(refreshNew)
	s := &entities.Session{
		ID:               uuid.New(),
		UserID:           usr.ID,
		TokenHash:        accessHash,
		RefreshTokenHash: refreshHash,
		ExpiresAt:        time.Now().Add(time.Hour),
		RefreshExpiresAt: time.Now().Add(24 * time.Hour),
		IsActive:         true,
		CreatedAt:        time.Now(),
	}
	u.sessions.Create(ctx, s)
	logger.Log.WithField("username", claims.Username).Info("session refreshed")
	return access, refreshNew, nil
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
	sess, err := u.sessions.GetByTokenHash(ctx, utils.HashString(token))
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
