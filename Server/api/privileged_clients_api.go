package api

import (
	context "context"
	"errors"
	"log"

	"api.mooody.me/db"
	"api.mooody.me/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkPrivilegedClient(ctx context.Context, clientUuid string, requirePrivileged bool) (*models.APIClient, error) {
	client, err := db.GetClientByUUID(ctx, clientUuid)

	if err != nil {
		log.Printf("database returns error '%s'.", err.Error())
		return nil, errors.New("internal error")
	}

	if !client.GetEnabled() {
		log.Printf("client '%s' is not enabled.", clientUuid)
		return nil, errors.New("client not enabled")
	}

	if requirePrivileged {
		if !client.GetPrivileged() {
			log.Printf("'requirePrivileged' was set, but client '%s' isn't privileged.", clientUuid)
			return nil, errors.New("unauthenticated")
		}
	}

	return client, nil
}

func (s *MoodyAPIServer) ListClients(ctx context.Context, request *models.ListClientsRequest) (*models.ListClientsResponse, error) {
	log.Printf("Received ListClients request.")
	_, err := checkPrivilegedClient(ctx, request.Auth.ClientId, true)

	if err != nil {
		log.Printf("checkPrivilegedClient failed: %s", err.Error())
		return &models.ListClientsResponse{Success: false}, errors.New("unauthenticated")
	}

	clients, err := db.ListClients(ctx)
	if err != nil {
		log.Printf("ListClients failed: %s", err.Error())
		return nil, errors.New("server error")
	}

	return &models.ListClientsResponse{Success: true, Clients: clients}, nil
}

func (s *MoodyAPIServer) UpdateClientInfo(ctx context.Context, request *models.UpdateClientInfoRequest) (*models.UpdateClientInfoResponse, error) {
	log.Printf("Received UpdateClientInfo request.")

	_, err := checkPrivilegedClient(ctx, request.Auth.ClientId, true)

	if err != nil {
		log.Printf("checkPrivilegedClient failed: %s", err.Error())
		return &models.UpdateClientInfoResponse{Success: false}, errors.New("unauthenticated")
	}

	client, err := db.GetClientByID(ctx, request.ClientInfo.Id)
	if err != nil {
		log.Printf("GetClientByID failed: %s", err.Error())
		return nil, errors.New("server error")
	}

	shouldReject := false

	if request.ClientInfo.Enabled != nil {
		shouldReject = shouldReject || request.Auth.ClientId == client.GetUuid()
		client.Enabled = request.ClientInfo.Enabled
	}

	if request.ClientInfo.LastSeen != nil {
		client.LastSeen = request.ClientInfo.LastSeen
	}

	if request.ClientInfo.Name != nil {
		client.Name = request.ClientInfo.Name
	}

	if request.ClientInfo.Privileged != nil {
		shouldReject = shouldReject || request.Auth.ClientId == client.GetUuid()
		client.Privileged = request.ClientInfo.Privileged
	}

	if request.ClientInfo.Uuid != nil {
		shouldReject = shouldReject || request.Auth.ClientId == client.GetUuid()
		client.Uuid = request.ClientInfo.Uuid
	}

	if shouldReject {
		log.Printf("client %s is performing suicide, reject", request.Auth.ClientId)
		return &models.UpdateClientInfoResponse{Success: false}, errors.New("don't suicide")
	}

	err = db.UpdateClient(ctx, client)

	if err != nil {
		log.Printf("UpdateClient failed: %s", err.Error())
		return &models.UpdateClientInfoResponse{Success: false}, errors.New("unexpected database result")
	}

	return &models.UpdateClientInfoResponse{Success: true}, nil
}

func (s *MoodyAPIServer) DeleteClient(context.Context, *models.DeleteClientRequest) (*models.DeleteClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteClient not implemented")
}
