package impl

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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
                       resource_name, ip_address, success, created_at)
               VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		a.ID, userID, a.Username, a.UserGroup, a.Action, a.ResourceType,
		a.ResourceName, a.IPAddress, a.Success, a.CreatedAt,
	)
	return err
}

func (r *pgAuditRepo) List(ctx context.Context, f *entities.AuditFilter) ([]*entities.AuditLog, int, error) {
	if f == nil {
		f = &entities.AuditFilter{}
	}
	f.SetDefaults()

	base := `SELECT id, user_id, username, user_group, action, resource_type,
                        resource_name, ip_address, success, created_at FROM audit_logs`
	countBase := `SELECT COUNT(1) FROM audit_logs`

	clauses := []string{}
	args := []interface{}{}
	idx := 1
	if f.Username != "" {
		clauses = append(clauses, "username=$"+strconv.Itoa(idx))
		args = append(args, f.Username)
		idx++
	}
	if f.UserGroup != "" {
		clauses = append(clauses, "user_group=$"+strconv.Itoa(idx))
		args = append(args, f.UserGroup)
		idx++
	}
	if f.IPAddress != "" {
		clauses = append(clauses, "ip_address=$"+strconv.Itoa(idx))
		args = append(args, f.IPAddress)
		idx++
	}
	if f.Resource != "" {
		clauses = append(clauses, "resource_type=$"+strconv.Itoa(idx))
		args = append(args, f.Resource)
		idx++
	}
	if f.FromTime != nil {
		clauses = append(clauses, "created_at >= $"+strconv.Itoa(idx))
		args = append(args, *f.FromTime)
		idx++
	}
	if f.ToTime != nil {
		clauses = append(clauses, "created_at <= $"+strconv.Itoa(idx))
		args = append(args, *f.ToTime)
		idx++
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	query := base + where + " ORDER BY created_at DESC" +
		fmt.Sprintf(" LIMIT %d OFFSET %d", f.Limit, f.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var logs []*entities.AuditLog
	for rows.Next() {
		var a entities.AuditLog
		if err := rows.Scan(&a.ID, &a.UserID, &a.Username, &a.UserGroup,
			&a.Action, &a.ResourceType, &a.ResourceName, &a.IPAddress,
			&a.Success, &a.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, &a)
	}

	countQuery := countBase + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *pgAuditRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.AuditLog, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, username, user_group, action, resource_type,
                        resource_name, ip_address, success, created_at
                FROM audit_logs WHERE id=$1`, id)
	var a entities.AuditLog
	err := row.Scan(
		&a.ID, &a.UserID, &a.Username, &a.UserGroup, &a.Action, &a.ResourceType,
		&a.ResourceName, &a.IPAddress, &a.Success, &a.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}
