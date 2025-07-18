package impl

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type pgPermissionRepo struct{ db *sql.DB }

func NewPermissionRepositoryPG(db *sql.DB) repositories.PermissionRepository {
	return &pgPermissionRepo{db: db}
}

func (r *pgPermissionRepo) List(ctx context.Context) ([]*entities.Permission, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, resource, action, description FROM permissions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []*entities.Permission
	for rows.Next() {
		var p entities.Permission
		if err := rows.Scan(&p.ID, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		perms = append(perms, &p)
	}
	return perms, nil
}

func (r *pgPermissionRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, resource, action, description FROM permissions WHERE id=$1`, id)
	var p entities.Permission
	if err := row.Scan(&p.ID, &p.Resource, &p.Action, &p.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgPermissionRepo) Create(ctx context.Context, p *entities.Permission) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO permissions (id, resource, action, description) VALUES ($1,$2,$3,$4)`,
		p.ID, p.Resource, p.Action, p.Description,
	)
	return err
}

func (r *pgPermissionRepo) Update(ctx context.Context, p *entities.Permission) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE permissions SET resource=$1, action=$2, description=$3 WHERE id=$4`,
		p.Resource, p.Action, p.Description, p.ID,
	)
	return err
}

func (r *pgPermissionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM permissions WHERE id=$1`, id)
	return err
}

func (r *pgPermissionRepo) GetByResourceAction(ctx context.Context, resource, action string) (*entities.Permission, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, resource, action, description FROM permissions WHERE resource=$1 AND action=$2`, resource, action)
	var p entities.Permission
	if err := row.Scan(&p.ID, &p.Resource, &p.Action, &p.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *pgPermissionRepo) GetByGroup(ctx context.Context, groupID uuid.UUID) ([]*entities.Permission, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT p.id, p.resource, p.action, p.description
        FROM permissions p
        JOIN group_permissions gp ON gp.permission_id = p.id
        WHERE gp.group_id=$1`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var perms []*entities.Permission
	for rows.Next() {
		var p entities.Permission
		if err := rows.Scan(&p.ID, &p.Resource, &p.Action, &p.Description); err != nil {
			return nil, err
		}
		perms = append(perms, &p)
	}
	return perms, nil
}

func (r *pgPermissionRepo) SetForGroup(ctx context.Context, groupID uuid.UUID, permIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM group_permissions WHERE group_id=$1`, groupID); err != nil {
		tx.Rollback()
		return err
	}
	for _, pid := range permIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO group_permissions (group_id, permission_id) VALUES ($1,$2)`, groupID, pid); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *pgPermissionRepo) HasGroupPermission(ctx context.Context, groupName, resource, action string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(1)
        FROM group_permissions gp
        JOIN groups g ON g.id = gp.group_id
        JOIN permissions p ON p.id = gp.permission_id
        WHERE g.name=$1 AND p.resource=$2 AND p.action=$3`, groupName, resource, action).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
