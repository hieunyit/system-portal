package impl

import (
	"context"
	"database/sql"

	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
)

type pgEmailTemplateRepo struct{ db *sql.DB }

func NewEmailTemplateRepositoryPG(db *sql.DB) repositories.EmailTemplateRepository {
	return &pgEmailTemplateRepo{db: db}
}

func (r *pgEmailTemplateRepo) GetByAction(ctx context.Context, action string) (*entities.EmailTemplate, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, action, subject, body, created_at, updated_at FROM email_templates WHERE action=$1`, action)
	var tpl entities.EmailTemplate
	if err := row.Scan(&tpl.ID, &tpl.Action, &tpl.Subject, &tpl.Body, &tpl.CreatedAt, &tpl.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tpl, nil
}

func (r *pgEmailTemplateRepo) Upsert(ctx context.Context, tpl *entities.EmailTemplate) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO email_templates (id, action, subject, body, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$5)
        ON CONFLICT (action) DO UPDATE SET subject=EXCLUDED.subject, body=EXCLUDED.body, updated_at=EXCLUDED.updated_at`,
		tpl.ID, tpl.Action, tpl.Subject, tpl.Body, tpl.UpdatedAt,
	)
	return err
}
