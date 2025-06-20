package repositories

import (
	"context"
	"system-portal/internal/domains/openvpn/entities"
)

type VPNStatusRepository interface {
	GetConnectedUsers(ctx context.Context) ([]*entities.ConnectedUser, error)
	IsUserConnected(ctx context.Context, username string) (*entities.ConnectedUser, bool, error)
	GetVPNStatus(ctx context.Context) (*entities.VPNStatusSummary, error)
}
