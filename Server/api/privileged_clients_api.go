package api

import (
	context "context"
	"errors"
	"log"
	"math/rand"
	"time"

	"api.mooody.me/db"
	"api.mooody.me/models"
)

func randate() time.Time {
	min := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func (s *MoodyAPIServer) ListClients(ctx context.Context, request *models.ListClientsRequest) (*models.ListClientsResponse, error) {
	log.Printf("Received ListClients request.")
	_, err := db.CheckClientValidity(ctx, request.Auth.ClientId, true)

	if err != nil {
		return &models.ListClientsResponse{Success: false}, errors.New("unauthenticated")
	}

	clients, err := db.ListClients(ctx)
	if err != nil {
		return nil, errors.New("server error")
	}

	return &models.ListClientsResponse{Success: true, Clients: clients}, nil
}

func (s *MoodyAPIServer) UpdateClientInfo(context.Context, *models.UpdateClientInfoRequest) (*models.UpdateClientInfoResponse, error) {
	log.Printf("Received UpdateClientInfo request.")
	return &models.UpdateClientInfoResponse{Success: true}, nil
}
