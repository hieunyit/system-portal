package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"system-portal/internal/domains/portal/entities"
	"system-portal/internal/domains/portal/repositories"
	"system-portal/internal/shared/infrastructure/ldap"
	"system-portal/internal/shared/infrastructure/xmlrpc"
)

type configUsecaseImpl struct {
	ovRepo       repositories.OpenVPNConfigRepository
	ldapRepo     repositories.LDAPConfigRepository
	smtpRepo     repositories.SMTPConfigRepository
	templateRepo repositories.EmailTemplateRepository
}

func NewConfigUsecase(ov repositories.OpenVPNConfigRepository, ldap repositories.LDAPConfigRepository, smtp repositories.SMTPConfigRepository, tpl repositories.EmailTemplateRepository) ConfigUsecase {
	return &configUsecaseImpl{ovRepo: ov, ldapRepo: ldap, smtpRepo: smtp, templateRepo: tpl}
}

func (u *configUsecaseImpl) GetOpenVPN(ctx context.Context) (*entities.OpenVPNConfig, error) {
	if u.ovRepo == nil {
		return nil, nil
	}
	return u.ovRepo.Get(ctx)
}

func (u *configUsecaseImpl) SetOpenVPN(ctx context.Context, cfg *entities.OpenVPNConfig) error {
	if u.ovRepo == nil {
		return nil
	}
	existing, err := u.ovRepo.Get(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	if existing == nil {
		if cfg.ID == uuid.Nil {
			cfg.ID = uuid.New()
		}
		cfg.CreatedAt = now
		cfg.UpdatedAt = now
		return u.ovRepo.Create(ctx, cfg)
	}
	cfg.ID = existing.ID
	cfg.CreatedAt = existing.CreatedAt
	cfg.UpdatedAt = now
	return u.ovRepo.Update(ctx, cfg)
}

func (u *configUsecaseImpl) DeleteOpenVPN(ctx context.Context) error {
	if u.ovRepo == nil {
		return nil
	}
	return u.ovRepo.Delete(ctx)
}

func (u *configUsecaseImpl) TestOpenVPN(ctx context.Context, cfg *entities.OpenVPNConfig) error {
	client := xmlrpc.NewClient(xmlrpc.Config{Host: cfg.Host, Username: cfg.Username, Password: cfg.Password, Port: cfg.Port})
	return client.Ping()
}

func (u *configUsecaseImpl) GetLDAP(ctx context.Context) (*entities.LDAPConfig, error) {
	if u.ldapRepo == nil {
		return nil, nil
	}
	return u.ldapRepo.Get(ctx)
}

func (u *configUsecaseImpl) SetLDAP(ctx context.Context, cfg *entities.LDAPConfig) error {
	if u.ldapRepo == nil {
		return nil
	}
	existing, err := u.ldapRepo.Get(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	if existing == nil {
		if cfg.ID == uuid.Nil {
			cfg.ID = uuid.New()
		}
		cfg.CreatedAt = now
		cfg.UpdatedAt = now
		return u.ldapRepo.Create(ctx, cfg)
	}
	cfg.ID = existing.ID
	cfg.CreatedAt = existing.CreatedAt
	cfg.UpdatedAt = now
	return u.ldapRepo.Update(ctx, cfg)
}

func (u *configUsecaseImpl) DeleteLDAP(ctx context.Context) error {
	if u.ldapRepo == nil {
		return nil
	}
	return u.ldapRepo.Delete(ctx)
}

func (u *configUsecaseImpl) TestLDAP(ctx context.Context, cfg *entities.LDAPConfig) error {
	client := ldap.NewClient(ldap.Config{Host: cfg.Host, Port: cfg.Port, BindDN: cfg.BindDN, BindPassword: cfg.BindPassword, BaseDN: cfg.BaseDN})
	conn, err := client.Connect()
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func (u *configUsecaseImpl) GetSMTP(ctx context.Context) (*entities.SMTPConfig, error) {
	if u.smtpRepo == nil {
		return nil, nil
	}
	return u.smtpRepo.Get(ctx)
}

func (u *configUsecaseImpl) SetSMTP(ctx context.Context, cfg *entities.SMTPConfig) error {
	if u.smtpRepo == nil {
		return nil
	}
	existing, err := u.smtpRepo.Get(ctx)
	if err != nil {
		return err
	}
	now := time.Now()
	if existing == nil {
		if cfg.ID == uuid.Nil {
			cfg.ID = uuid.New()
		}
		cfg.CreatedAt = now
		cfg.UpdatedAt = now
		return u.smtpRepo.Create(ctx, cfg)
	}
	cfg.ID = existing.ID
	cfg.CreatedAt = existing.CreatedAt
	cfg.UpdatedAt = now
	return u.smtpRepo.Update(ctx, cfg)
}

func (u *configUsecaseImpl) DeleteSMTP(ctx context.Context) error {
	if u.smtpRepo == nil {
		return nil
	}
	return u.smtpRepo.Delete(ctx)
}

func (u *configUsecaseImpl) GetTemplate(ctx context.Context, action string) (*entities.EmailTemplate, error) {
	if u.templateRepo == nil {
		return nil, nil
	}
	return u.templateRepo.GetByAction(ctx, action)
}

func (u *configUsecaseImpl) SetTemplate(ctx context.Context, tpl *entities.EmailTemplate) error {
	if u.templateRepo == nil {
		return nil
	}
	if tpl.ID == uuid.Nil {
		tpl.ID = uuid.New()
	}
	tpl.UpdatedAt = time.Now()
	if tpl.CreatedAt.IsZero() {
		tpl.CreatedAt = tpl.UpdatedAt
	}
	return u.templateRepo.Upsert(ctx, tpl)
}
