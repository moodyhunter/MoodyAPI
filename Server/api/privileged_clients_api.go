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
	_, err := checkPrivilegedClient(ctx, request.Auth.ClientUuid, true)

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

func (s *MoodyAPIServer) UpdateClient(ctx context.Context, request *models.UpdateClientRequest) (*models.UpdateClientResponse, error) {
	_, err := checkPrivilegedClient(ctx, request.Auth.ClientUuid, true)

	if err != nil {
		log.Printf("checkPrivilegedClient failed: %s", err.Error())
		return &models.UpdateClientResponse{Success: false}, errors.New("unauthenticated")
	}

	client, err := db.GetClientByID(ctx, request.Client.Id)
	if err != nil {
		log.Printf("GetClientByID failed: %s", err.Error())
		return nil, errors.New("server error")
	}

	shouldReject := false

	if request.Client.Enabled != nil {
		shouldReject = shouldReject || request.Auth.ClientUuid == client.GetUuid()
		client.Enabled = request.Client.Enabled
	}

	if request.Client.LastSeen != nil {
		client.LastSeen = request.Client.LastSeen
	}

	if request.Client.Name != nil {
		client.Name = request.Client.Name
	}

	if request.Client.Privileged != nil {
		shouldReject = shouldReject || request.Auth.ClientUuid == client.GetUuid()
		client.Privileged = request.Client.Privileged
	}

	if request.Client.Uuid != nil {
		shouldReject = shouldReject || request.Auth.ClientUuid == client.GetUuid()
		client.Uuid = request.Client.Uuid
	}

	if shouldReject {
		log.Printf("client %s is performing suicide, reject", request.Auth.ClientUuid)
		return &models.UpdateClientResponse{Success: false}, errors.New("don't suicide")
	}

	err = db.UpdateClient(ctx, client)

	if err != nil {
		log.Printf("UpdateClient failed: %s", err.Error())
		return &models.UpdateClientResponse{Success: false}, errors.New("unexpected database result")
	}

	return &models.UpdateClientResponse{Success: true}, nil
}

func (s *MoodyAPIServer) DeleteClient(ctx context.Context, request *models.DeleteClientRequest) (*models.DeleteClientResponse, error) {
	_, err := checkPrivilegedClient(ctx, request.Auth.ClientUuid, true)

	if err != nil {
		log.Printf("checkPrivilegedClient failed: %s", err.Error())
		return &models.DeleteClientResponse{Success: false}, errors.New("unauthenticated")
	}

	client, err := db.GetClientByID(ctx, request.Client.Id)
	if err != nil {
		log.Printf("GetClientByID failed: %s", err.Error())
		return nil, errors.New("server error")
	}

	if request.Auth.ClientUuid == *client.Uuid {
		log.Printf("client %s is performing suicide, reject", request.Auth.ClientUuid)
		return &models.DeleteClientResponse{Success: false}, errors.New("don't suicide")
	}

	err = db.DeleteClient(ctx, client)

	if err != nil {
		return &models.DeleteClientResponse{Success: false}, errors.New("unexpected database result")
	}
	return &models.DeleteClientResponse{Success: true}, nil
}
