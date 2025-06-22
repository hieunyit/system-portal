package impl

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type pgAuditRepo struct{ db *sql.DB }

func NewAuditRepositoryPG(db *sql.DB) repositories.AuditRepository {
	return &pgAuditRepo{db: db}
}

func (r *pgAuditRepo) Add(ctx context.Context, a *entities.AuditLog) error {
	var userID interface{}
	if a.UserID == uuid.Nil {
		userID = nil
	} else {
		userID = a.UserID
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO audit_logs (
                        id, user_id, username, user_group, action, resource_type,
                        resource_id, resource_name, ip_address, user_agent,
                        success, error_message, duration_ms, created_at)
                VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		a.ID, userID, a.Username, a.UserGroup, a.Action, a.Resource,
		a.ResourceID, a.ResourceName, a.IPAddress, a.UserAgent,
		a.Success, a.ErrorMessage, a.DurationMs, a.CreatedAt,
	)
	return err
}

func (r *pgAuditRepo) List(ctx context.Context) ([]*entities.AuditLog, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, username, user_group, action, resource_type,
                        resource_id, resource_name, ip_address, user_agent,
                        success, error_message, duration_ms, created_at
                FROM audit_logs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []*entities.AuditLog
	for rows.Next() {
		var a entities.AuditLog
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Username, &a.UserGroup, &a.Action, &a.Resource,
			&a.ResourceID, &a.ResourceName, &a.IPAddress, &a.UserAgent,
			&a.Success, &a.ErrorMessage, &a.DurationMs, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, &a)
	}
	return logs, nil
}

func (r *pgAuditRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, username, user_group, action, resource_type,
                        resource_id, resource_name, ip_address, user_agent,
                        success, error_message, duration_ms, created_at
                FROM audit_logs WHERE id=$1`, id)
	var a entities.AuditLog
	err := row.Scan(
		&a.ID, &a.UserID, &a.Username, &a.UserGroup, &a.Action, &a.Resource,
		&a.ResourceID, &a.ResourceName, &a.IPAddress, &a.UserAgent,
		&a.Success, &a.ErrorMessage, &a.DurationMs, &a.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
