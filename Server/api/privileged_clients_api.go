package api

import (
	context "context"
	"errors"
	"log"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MoodyAPIServer) ListClients(ctx context.Context, request *models.ListClientsRequest) (*models.ListClientsResponse, error) {
	client, err := common.GetClientFromAuth(ctx, request.Auth, true)
	if err != nil {
		common.LogClientWithError(client, err)
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
	client, err := common.GetClientFromAuth(ctx, request.Auth, true)
	if err != nil {
		common.LogClientWithError(client, err)
		return &models.UpdateClientResponse{Success: false}, errors.New("unauthenticated")
	}

	if request.Client == nil {
		return nil, errors.New("invalid request")
	}

	client, err = db.GetClientByID(ctx, request.Client.Id)
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
		log.Printf("client %s is performing suicide, reject", *client.Name)
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
	client, err := common.GetClientFromAuth(ctx, request.Auth, true)
	if err != nil {

		common.LogClientWithError(client, err)
		return &models.DeleteClientResponse{Success: false}, errors.New("unauthenticated")
	}

	if request.Client == nil {
		return nil, errors.New("invalid request")
	}

	client, err = db.GetClientByID(ctx, request.Client.Id)
	if err != nil {
		log.Printf("GetClientByID failed: %s", err.Error())
		return nil, errors.New("server error")
	}

	if request.Auth.ClientUuid == *client.Uuid {
		log.Printf("client %s is performing suicide, reject", *client.Name)
		return &models.DeleteClientResponse{Success: false}, errors.New("don't suicide")
	}

	err = db.DeleteClient(ctx, client)

	if err != nil {
		log.Printf("deleteClient failed: %s", err)
		return &models.DeleteClientResponse{Success: false}, errors.New("unexpected database result")
	}
	return &models.DeleteClientResponse{Success: true}, nil
}

func (s *MoodyAPIServer) CreateClient(ctx context.Context, request *models.CreateClientRequest) (*models.CreateClientResponse, error) {
	_, err := common.GetClientFromAuth(ctx, request.Auth, true)
	if err != nil {
		log.Printf("checkPrivilegedClient failed: %s", err.Error())
		return &models.CreateClientResponse{Success: false}, errors.New("unauthenticated")
	}

	if request.Client == nil {
		return nil, errors.New("invalid request")
	}

	request.Client.Privileged = proto.Bool(false)
	request.Client.Enabled = proto.Bool(true)
	request.Client.Uuid = proto.String(uuid.New().String())
	request.Client.Id = 0
	request.Client.LastSeen = timestamppb.Now()

	client, err := db.CreateClient(ctx, request.Client)
	if err != nil {
		log.Printf("createClient failed: %s", err)
		return &models.CreateClientResponse{Success: false}, errors.New("unexpected database result")
	}

	return &models.CreateClientResponse{Success: true, Client: client}, nil
}
