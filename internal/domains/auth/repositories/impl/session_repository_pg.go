package impl

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"system-portal/internal/domains/auth/entities"
	"system-portal/internal/domains/auth/repositories"
	"system-portal/pkg/logger"
)

type pgSessionRepo struct{ db *sql.DB }

func NewSessionRepositoryPG(db *sql.DB) repositories.SessionRepository {
	return &pgSessionRepo{db: db}
}

func (r *pgSessionRepo) Create(ctx context.Context, s *entities.Session) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user_sessions (id, user_id, token_hash, refresh_token_hash, expires_at, refresh_expires_at, is_active, created_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		s.ID, s.UserID, s.TokenHash, s.RefreshTokenHash, s.ExpiresAt, s.RefreshExpiresAt, s.IsActive, s.CreatedAt,
	)
	if err != nil {
		logger.Log.WithError(err).Error("create session failed")
	}
	return err
}

func (r *pgSessionRepo) GetByTokenHash(ctx context.Context, hash string) (*entities.Session, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token_hash, refresh_token_hash, expires_at, refresh_expires_at, is_active, created_at
         FROM user_sessions WHERE token_hash=$1`, hash)
	var s entities.Session
	err := row.Scan(&s.ID, &s.UserID, &s.TokenHash, &s.RefreshTokenHash, &s.ExpiresAt, &s.RefreshExpiresAt, &s.IsActive, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Log.WithError(err).Error("get session failed")
		return nil, err
	}
	return &s, nil
}

func (r *pgSessionRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE user_sessions SET is_active=false WHERE id=$1`, id)
	if err != nil {
		logger.Log.WithError(err).Error("deactivate session failed")
	}
	return err
}
