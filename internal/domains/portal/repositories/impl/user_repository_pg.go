package impl

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
	"system-portal/pkg/logger"

	"github.com/google/uuid"
)

type pgUserRepo struct{ db *sql.DB }

func NewUserRepositoryPG(db *sql.DB) repositories.UserRepository {
	return &pgUserRepo{db: db}
}

func (r *pgUserRepo) Create(ctx context.Context, u *entities.PortalUser) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, username, email, password_hash, full_name, group_id, is_active, created_at, updated_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		u.ID, u.Username, u.Email, u.Password, u.FullName, u.GroupID, u.IsActive, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		logger.Log.WithError(err).Error("create user failed")
	}
	return err
}

func (r *pgUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.PortalUser, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, full_name, group_id, is_active, created_at, updated_at
        FROM users WHERE id=$1`, id)
	var u entities.PortalUser
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.GroupID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Log.WithError(err).Error("get user by id failed")
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepo) GetByUsername(ctx context.Context, username string) (*entities.PortalUser, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, full_name, group_id, is_active, created_at, updated_at
        FROM users WHERE username=$1`, username)
	var u entities.PortalUser
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.FullName, &u.GroupID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Log.WithError(err).Error("get user by username failed")
		return nil, err
	}
	logger.Log.WithField("username", u.Username).Debug("user retrieved")
	return &u, nil
}

func (r *pgUserRepo) GetByEmail(ctx context.Context, email string) (*entities.PortalUser, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, full_name, group_id, is_active, created_at, updated_at
        FROM users WHERE email=$1`, email)
	var u entities.PortalUser
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.FullName, &u.GroupID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Log.WithError(err).Error("get user by email failed")
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepo) List(ctx context.Context, f *entities.UserFilter) ([]*entities.PortalUser, int, error) {
	if f == nil {
		f = &entities.UserFilter{}
	}
	f.SetDefaults()

	base := `SELECT id, username, email, full_name, group_id, is_active, created_at, updated_at FROM users`
	countBase := `SELECT COUNT(1) FROM users`
	clauses := []string{}
	args := []interface{}{}
	idx := 1
	if f.Username != "" {
		clauses = append(clauses, "username ILIKE $"+strconv.Itoa(idx))
		args = append(args, "%"+f.Username+"%")
		idx++
	}
	if f.Email != "" {
		clauses = append(clauses, "email ILIKE $"+strconv.Itoa(idx))
		args = append(args, "%"+f.Email+"%")
		idx++
	}
	if f.GroupID != uuid.Nil {
		clauses = append(clauses, "group_id=$"+strconv.Itoa(idx))
		args = append(args, f.GroupID)
		idx++
	}
	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}
	query := base + where + fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", f.Limit, f.Offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var users []*entities.PortalUser
	for rows.Next() {
		var u entities.PortalUser
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.GroupID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, &u)
	}
	countQuery := countBase + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *pgUserRepo) Update(ctx context.Context, u *entities.PortalUser) error {
	u.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash=$2, full_name=$3, group_id=$4, is_active=$5, updated_at=$6 WHERE id=$1`,
		u.ID, u.Password, u.FullName, u.GroupID, u.IsActive, u.UpdatedAt,
	)
	if err != nil {
		logger.Log.WithError(err).Error("update user failed")
	}
	return err
}

func (r *pgUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		logger.Log.WithError(err).Error("delete user failed")
	}
	return err
}
