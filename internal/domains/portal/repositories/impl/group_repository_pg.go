package impl

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type pgGroupRepo struct{ db *sql.DB }

func NewGroupRepositoryPG(db *sql.DB) repositories.GroupRepository {
	return &pgGroupRepo{db: db}
}

func (r *pgGroupRepo) Create(ctx context.Context, g *entities.PortalGroup) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO groups (id, name, display_name, is_active, created_at, updated_at)
         VALUES ($1,$2,$3,$4,$5,$6)`,
		g.ID, g.Name, g.DisplayName, g.IsActive, g.CreatedAt, g.UpdatedAt,
	)
	return err
}

func (r *pgGroupRepo) GetByID(ctx context.Context, id uuid.UUID) (*entities.PortalGroup, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, display_name, is_active, created_at, updated_at FROM groups WHERE id=$1`, id)
	var g entities.PortalGroup
	err := row.Scan(&g.ID, &g.Name, &g.DisplayName, &g.IsActive, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *pgGroupRepo) GetByName(ctx context.Context, name string) (*entities.PortalGroup, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, display_name, is_active, created_at, updated_at FROM groups WHERE name=$1`, name)
	var g entities.PortalGroup
	err := row.Scan(&g.ID, &g.Name, &g.DisplayName, &g.IsActive, &g.CreatedAt, &g.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *pgGroupRepo) List(ctx context.Context) ([]*entities.PortalGroup, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, display_name, is_active, created_at, updated_at FROM groups`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []*entities.PortalGroup
	for rows.Next() {
		var g entities.PortalGroup
		if err := rows.Scan(&g.ID, &g.Name, &g.DisplayName, &g.IsActive, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}
	return groups, nil
}
