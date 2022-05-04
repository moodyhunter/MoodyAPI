package api

import (
	"context"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models/wg"
)

// CreateWireguardClient creates a new wireguard client, implements the MoodyAPIServer interface
func (s *MoodyAPIServer) CreateWireguardClient(ctx context.Context, request *wg.CreateClientRequest) (*wg.CreateClientResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `creates wireguard client`)

	// TODO
	return nil, nil
}

// DeleteWireguardClient
func (s *MoodyAPIServer) DeleteWireguardClient(ctx context.Context, request *wg.DeleteClientRequest) (*wg.DeleteClientResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `deletes wireguard client`)

	// TODO
	return nil, nil
}

// ListWireguardClients
func (s *MoodyAPIServer) ListWireguardClients(ctx context.Context, request *wg.ListWireguardClientsRequest) (*wg.ListWireguardClientsResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `lists wireguard clients`)

	// TODO
	return nil, nil
}

// UpdateWireguardClient
func (s *MoodyAPIServer) UpdateWireguardClient(ctx context.Context, request *wg.UpdateClientRequest) (*wg.UpdateClientResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `updates wireguard client`)

	// TODO
	return nil, nil
}
