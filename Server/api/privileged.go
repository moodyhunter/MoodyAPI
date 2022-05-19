package api

import (
	"context"
	"errors"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models/privileged"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *MoodyAPIServer) ListClients(ctx context.Context, request *privileged.ListClientsRequest) (*privileged.ListClientsResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, true)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return &privileged.ListClientsResponse{Success: false}, errors.New("unauthenticated")
	}

	clients, err := db.ListClients(ctx)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, errors.New("server error")
	}

	return &privileged.ListClientsResponse{Success: true, Clients: clients}, nil
}

func (s *MoodyAPIServer) UpdateClient(ctx context.Context, request *privileged.UpdateClientRequest) (*privileged.UpdateClientResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, true)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return &privileged.UpdateClientResponse{Success: false}, errors.New("unauthenticated")
	}

	if request.Client == nil {
		return nil, errors.New("invalid request")
	}

	client, err = db.GetClientByID(ctx, request.Client.Id)
	if err != nil {
		common.LogClientError(ctx, client, err)
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
		common.LogClientOperation(ctx, client, "client is performing suicide, rejecting.")
		return &privileged.UpdateClientResponse{Success: false}, errors.New("don't suicide")
	}

	err = db.UpdateClient(ctx, client)

	if err != nil {
		common.LogClientError(ctx, client, err)
		return &privileged.UpdateClientResponse{Success: false}, errors.New("unexpected database result")
	}

	common.LogClientOperation(ctx, client, "updated client with id '%d'.", request.Client.Id)
	return &privileged.UpdateClientResponse{Success: true}, nil
}

func (s *MoodyAPIServer) DeleteClient(ctx context.Context, request *privileged.DeleteClientRequest) (*privileged.DeleteClientResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, true)
	if err != nil {

		common.LogClientError(ctx, client, err)
		return &privileged.DeleteClientResponse{Success: false}, errors.New("unauthenticated")
	}

	if request.Client == nil {
		return nil, errors.New("invalid request")
	}

	client, err = db.GetClientByID(ctx, request.Client.Id)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, errors.New("server error")
	}

	if request.Auth.ClientUuid == *client.Uuid {
		common.LogClientOperation(ctx, client, "client is performing suicide, rejecting.")
		return &privileged.DeleteClientResponse{Success: false}, errors.New("don't suicide")
	}

	err = db.DeleteClient(ctx, client)

	if err != nil {
		common.LogClientError(ctx, client, err)
		return &privileged.DeleteClientResponse{Success: false}, errors.New("unexpected database result")
	}

	common.LogClientOperation(ctx, client, "deleted client with id '%d'.", request.Client.Id)
	return &privileged.DeleteClientResponse{Success: true}, nil
}

func (s *MoodyAPIServer) CreateClient(ctx context.Context, request *privileged.CreateClientRequest) (*privileged.CreateClientResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, true)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return &privileged.CreateClientResponse{Success: false}, errors.New("unauthenticated")
	}

	if request.Client == nil {
		return nil, errors.New("invalid request")
	}

	request.Client.Privileged = proto.Bool(false)
	request.Client.Enabled = proto.Bool(true)
	request.Client.Uuid = proto.String(uuid.New().String())
	request.Client.Id = 0
	request.Client.LastSeen = timestamppb.Now()

	newClient, err := db.CreateClient(ctx, request.Client)
	if err != nil {
		common.LogClientError(ctx, newClient, err)
		return &privileged.CreateClientResponse{Success: false}, errors.New("unexpected database result")
	}

	common.LogClientOperation(ctx, client, "created another client with id '%d', named '%s'.", newClient.Id, newClient.GetName())
	return &privileged.CreateClientResponse{Success: true, Client: newClient}, nil
}
