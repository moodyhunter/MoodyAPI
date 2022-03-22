package api

import (
	context "context"
	"errors"
	"log"

	"api.mooody.me/db"
	"api.mooody.me/models"
)

func checkPrivilegedClient(ctx context.Context, clientUuid string, requirePrivileged bool) (*models.APIClient, error) {
	client, err := db.GetClientByUUID(ctx, clientUuid)

	if err != nil {
		return nil, errors.New("unauthenticated")
	}

	if !client.GetEnabled() {
		return nil, errors.New("client not enabled")
	}

	if requirePrivileged {
		if !client.GetPrivileged() {
			return nil, errors.New("unauthenticated")
		}
	}

	return client, nil
}

func (s *MoodyAPIServer) ListClients(ctx context.Context, request *models.ListClientsRequest) (*models.ListClientsResponse, error) {
	log.Printf("Received ListClients request.")
	_, err := checkPrivilegedClient(ctx, request.Auth.ClientId, true)

	if err != nil {
		return &models.ListClientsResponse{Success: false}, errors.New("unauthenticated")
	}

	clients, err := db.ListClients(ctx)
	if err != nil {
		return nil, errors.New("server error")
	}

	return &models.ListClientsResponse{Success: true, Clients: clients}, nil
}

func (s *MoodyAPIServer) UpdateClientInfo(ctx context.Context, request *models.UpdateClientInfoRequest) (*models.UpdateClientInfoResponse, error) {
	log.Printf("Received UpdateClientInfo request.")

	_, err := checkPrivilegedClient(ctx, request.Auth.ClientId, true)

	if err != nil {
		return &models.UpdateClientInfoResponse{Success: false}, errors.New("unauthenticated")
	}

	client, err := db.GetClientByID(ctx, request.ClientInfo.Id)
	if err != nil {
		return nil, errors.New("server error")
	}

	if request.ClientInfo.Enabled != nil {
		client.Enabled = request.ClientInfo.Enabled
	}

	if request.ClientInfo.LastSeen != nil {
		client.LastSeen = request.ClientInfo.LastSeen
	}

	if request.ClientInfo.Name != nil {
		client.Name = request.ClientInfo.Name
	}

	if request.ClientInfo.Privileged != nil {
		client.Privileged = request.ClientInfo.Privileged
	}

	if request.ClientInfo.Uuid != nil {
		client.Uuid = request.ClientInfo.Uuid
	}

	err = db.UpdateClient(ctx, client)

	if err != nil {
		return &models.UpdateClientInfoResponse{Success: false}, errors.New("unexpected database result")
	}

	return &models.UpdateClientInfoResponse{Success: true}, nil
}
