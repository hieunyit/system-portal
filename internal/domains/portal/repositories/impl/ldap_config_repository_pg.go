package impl

import (
	"context"
	"database/sql"

	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
	"system-portal/pkg/utils"
)

type pgLDAPConfigRepo struct {
	db  *sql.DB
	key string
}

func NewLDAPConfigRepositoryPG(db *sql.DB, key string) repositories.LDAPConfigRepository {
	return &pgLDAPConfigRepo{db: db, key: key}
}

func (r *pgLDAPConfigRepo) Get(ctx context.Context) (*entities.LDAPConfig, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, host, port, bind_dn, bind_password, base_dn, created_at, updated_at FROM ldap_configs LIMIT 1`)
	var cfg entities.LDAPConfig
	var encPass string
	if err := row.Scan(&cfg.ID, &cfg.Host, &cfg.Port, &cfg.BindDN, &encPass, &cfg.BaseDN, &cfg.CreatedAt, &cfg.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if r.key != "" {
		if plain, err := utils.DecryptString(encPass, r.key); err == nil {
			cfg.BindPassword = plain
		} else {
			cfg.BindPassword = encPass
		}
	} else {
		cfg.BindPassword = encPass
	}
	return &cfg, nil
}

func (r *pgLDAPConfigRepo) Create(ctx context.Context, cfg *entities.LDAPConfig) error {
	pass := cfg.BindPassword
	if r.key != "" {
		if enc, err := utils.EncryptString(cfg.BindPassword, r.key); err == nil {
			pass = enc
		} else {
			return err
		}
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ldap_configs (id, host, port, bind_dn, bind_password, base_dn, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$7)`,
		cfg.ID, cfg.Host, cfg.Port, cfg.BindDN, pass, cfg.BaseDN, cfg.CreatedAt,
	)
	return err
}

func (r *pgLDAPConfigRepo) Update(ctx context.Context, cfg *entities.LDAPConfig) error {
	pass := cfg.BindPassword
	if r.key != "" {
		if enc, err := utils.EncryptString(cfg.BindPassword, r.key); err == nil {
			pass = enc
		} else {
			return err
		}
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE ldap_configs SET host=$1, port=$2, bind_dn=$3, bind_password=$4, base_dn=$5, updated_at=$6 WHERE id=$7`,
		cfg.Host, cfg.Port, cfg.BindDN, pass, cfg.BaseDN, cfg.UpdatedAt, cfg.ID,
	)
	return err
}

func (r *pgLDAPConfigRepo) Delete(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ldap_configs`)
	return err
}
