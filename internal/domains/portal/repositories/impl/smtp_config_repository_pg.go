package impl

import (
	"context"
	"database/sql"

	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
	"system-portal/pkg/utils"
)

type pgSMTPConfigRepo struct {
	db  *sql.DB
	key string
}

func NewSMTPConfigRepositoryPG(db *sql.DB, key string) repositories.SMTPConfigRepository {
	return &pgSMTPConfigRepo{db: db, key: key}
}

func (r *pgSMTPConfigRepo) Get(ctx context.Context) (*entities.SMTPConfig, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, host, port, username, password, from_addr, tls, created_at, updated_at FROM smtp_configs LIMIT 1`)
	var cfg entities.SMTPConfig
	var encPass string
	if err := row.Scan(&cfg.ID, &cfg.Host, &cfg.Port, &cfg.Username, &encPass, &cfg.From, &cfg.TLS, &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
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

func (r *pgSMTPConfigRepo) Create(ctx context.Context, cfg *entities.SMTPConfig) error {
	pass := cfg.Password
	if r.key != "" {
		if enc, err := utils.EncryptString(cfg.Password, r.key); err == nil {
			pass = enc
		} else {
			return err
		}
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO smtp_configs (id, host, port, username, password, from_addr, tls, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8)`,
		cfg.ID, cfg.Host, cfg.Port, cfg.Username, pass, cfg.From, cfg.TLS, cfg.CreatedAt,
	)
	return err
}

func (r *pgSMTPConfigRepo) Update(ctx context.Context, cfg *entities.SMTPConfig) error {
	pass := cfg.Password
	if r.key != "" {
		if enc, err := utils.EncryptString(cfg.Password, r.key); err == nil {
			pass = enc
		} else {
			return err
		}
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE smtp_configs SET host=$1, port=$2, username=$3, password=$4, from_addr=$5, tls=$6, updated_at=$7 WHERE id=$8`,
		cfg.Host, cfg.Port, cfg.Username, pass, cfg.From, cfg.TLS, cfg.UpdatedAt, cfg.ID,
	)
	return err
}

func (r *pgSMTPConfigRepo) Delete(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM smtp_configs`)
	return err
}
