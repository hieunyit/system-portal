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
	ovRepo   repositories.OpenVPNConfigRepository
	ldapRepo repositories.LDAPConfigRepository
}

func NewConfigUsecase(ov repositories.OpenVPNConfigRepository, ldap repositories.LDAPConfigRepository) ConfigUsecase {
	return &configUsecaseImpl{ovRepo: ov, ldapRepo: ldap}
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
