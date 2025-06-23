package impl

import (
	"context"
	"database/sql"

	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
	"system-portal/pkg/utils"
)

type pgOpenVPNConfigRepo struct {
	db  *sql.DB
	key string
}

func NewOpenVPNConfigRepositoryPG(db *sql.DB, key string) repositories.OpenVPNConfigRepository {
	return &pgOpenVPNConfigRepo{db: db, key: key}
}

func (r *pgOpenVPNConfigRepo) Get(ctx context.Context) (*entities.OpenVPNConfig, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, host, username, password, port, created_at, updated_at FROM openvpn_configs LIMIT 1`)
	var cfg entities.OpenVPNConfig
	var encPass string
	if err := row.Scan(&cfg.ID, &cfg.Host, &cfg.Username, &encPass, &cfg.Port, &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if r.key != "" {
		if plain, err := utils.DecryptString(encPass, r.key); err == nil {
			cfg.Password = plain
		} else {
			cfg.Password = encPass
		}
	} else {
		cfg.Password = encPass
	}
	return &cfg, nil
}

func (r *pgOpenVPNConfigRepo) Create(ctx context.Context, cfg *entities.OpenVPNConfig) error {
	pass := cfg.Password
	if r.key != "" {
		if enc, err := utils.EncryptString(cfg.Password, r.key); err == nil {
			pass = enc
		} else {
			return err
		}
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO openvpn_configs (id, host, username, password, port, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$6)`,
		cfg.ID, cfg.Host, cfg.Username, pass, cfg.Port, cfg.CreatedAt,
	)
	return err
}

func (r *pgOpenVPNConfigRepo) Update(ctx context.Context, cfg *entities.OpenVPNConfig) error {
	pass := cfg.Password
	if r.key != "" {
		if enc, err := utils.EncryptString(cfg.Password, r.key); err == nil {
			pass = enc
		} else {
			return err
		}
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE openvpn_configs SET host=$1, username=$2, password=$3, port=$4, updated_at=$5 WHERE id=$6`,
		cfg.Host, cfg.Username, pass, cfg.Port, cfg.UpdatedAt, cfg.ID,
	)
	return err
}

func (r *pgOpenVPNConfigRepo) Delete(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM openvpn_configs`)
	return err
}
