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

func (r *pgGroupRepo) List(ctx context.Context, f *entities.GroupFilter) ([]*entities.PortalGroup, int, error) {
	if f == nil {
		f = &entities.GroupFilter{}
	}
	f.SetDefaults()

	base := `SELECT id, name, display_name, is_active, created_at, updated_at FROM groups`
	countBase := `SELECT COUNT(1) FROM groups`
	clauses := []string{}
	args := []interface{}{}
	idx := 1
	if f.Name != "" {
		clauses = append(clauses, "name ILIKE $"+strconv.Itoa(idx))
		args = append(args, "%"+f.Name+"%")
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
	var groups []*entities.PortalGroup
	for rows.Next() {
		var g entities.PortalGroup
		if err := rows.Scan(&g.ID, &g.Name, &g.DisplayName, &g.IsActive, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, 0, err
		}
		groups = append(groups, &g)
	}
	countQuery := countBase + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

func (r *pgGroupRepo) Update(ctx context.Context, g *entities.PortalGroup) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE groups SET name=$2, display_name=$3, is_active=$4, updated_at=$5 WHERE id=$1`,
		g.ID, g.Name, g.DisplayName, g.IsActive, g.UpdatedAt,
	)
	return err
}

func (r *pgGroupRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM groups WHERE id=$1`, id)
	return err
}
