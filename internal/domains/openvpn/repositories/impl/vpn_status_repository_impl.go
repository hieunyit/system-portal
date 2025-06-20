package repositories

import (
	"context"
	"strings"
	"system-portal/internal/domains/openvpn/entities"
	"system-portal/internal/domains/openvpn/repositories"
	"system-portal/internal/shared/infrastructure/xmlrpc"
)

type vpnStatusRepositoryImpl struct {
	vpnStatusClient *xmlrpc.VPNStatusClient
}

func NewVPNStatusRepository(xmlrpcClient *xmlrpc.Client) repositories.VPNStatusRepository {
	return &vpnStatusRepositoryImpl{
		vpnStatusClient: xmlrpc.NewVPNStatusClient(xmlrpcClient),
	}
}

func (r *vpnStatusRepositoryImpl) GetConnectedUsers(ctx context.Context) ([]*entities.ConnectedUser, error) {
	status, err := r.vpnStatusClient.GetVPNStatus()
	if err != nil {
		return nil, err
	}
	return status.ConnectedUsers, nil
}

func (r *vpnStatusRepositoryImpl) IsUserConnected(ctx context.Context, username string) (*entities.ConnectedUser, bool, error) {
	status, err := r.vpnStatusClient.GetVPNStatus()
	if err != nil {
		return nil, false, err
	}

	// Find user in connected users list (case-insensitive)
	for _, user := range status.ConnectedUsers {
		if strings.EqualFold(user.Username, username) {
			return user, true, nil
		}
	}

	return nil, false, nil
}

func (r *vpnStatusRepositoryImpl) GetVPNStatus(ctx context.Context) (*entities.VPNStatusSummary, error) {
	return r.vpnStatusClient.GetVPNStatus()
}
