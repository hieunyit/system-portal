package email

import (
	"bytes"
	"context"
	"text/template"

	portalRepo "system-portal/internal/domains/portal/repositories"
	"system-portal/internal/shared/config"
	"system-portal/pkg/mailer"
)

// Service handles sending emails using templates stored in the database
type Service struct {
	smtpRepo portalRepo.SMTPConfigRepository
	tplRepo  portalRepo.EmailTemplateRepository
}

func NewService(smtpRepo portalRepo.SMTPConfigRepository, tplRepo portalRepo.EmailTemplateRepository) *Service {
	return &Service{smtpRepo: smtpRepo, tplRepo: tplRepo}
}

func (s *Service) Send(ctx context.Context, action, to string, data map[string]interface{}) error {
	if s.smtpRepo == nil || s.tplRepo == nil {
		return nil
	}
	cfg, err := s.smtpRepo.Get(ctx)
	if err != nil || cfg == nil {
		return err
	}
	tpl, err := s.tplRepo.GetByAction(ctx, action)
	if err != nil || tpl == nil {
		return err
	}
	subject, err := executeTemplate(tpl.Subject, data)
	if err != nil {
		return err
	}
	body, err := executeTemplate(tpl.Body, data)
	if err != nil {
		return err
	}
	m := mailer.NewSMTPMailer(config.SMTPConfig{
		Host: cfg.Host, Port: cfg.Port, Username: cfg.Username, Password: cfg.Password, From: cfg.From, TLS: cfg.TLS,
	})
	return m.Send(to, subject, body)
}

func executeTemplate(tmpl string, data map[string]interface{}) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
