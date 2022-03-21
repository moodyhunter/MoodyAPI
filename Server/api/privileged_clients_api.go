package api

import (
	context "context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"api.mooody.me/api/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func randate() time.Time {
	min := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2022, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func (s *MoodyAPIServer) ListClients(context.Context, *pb.ListClientsRequest) (*pb.ListClientsResponse, error) {
	log.Printf("Received ListClients request.")

	clients := [20]*pb.APIClient{}

	for i := 0; i < 20; i++ {
		clients[i] = &pb.APIClient{
			Id:         int64(i),
			Name:       "Test Client: " + fmt.Sprint(i),
			Uuid:       uuid.New().String(),
			Privileged: false,
			LastSeen:   timestamppb.New(randate()),
			Enabled:    i%3 == 0,
		}
	}

	return &pb.ListClientsResponse{Success: true, Clients: clients[:]}, nil
}

func (s *MoodyAPIServer) UpdateClientInfo(context.Context, *pb.UpdateClientInfoRequest) (*pb.UpdateClientInfoResponse, error) {
	log.Printf("Received UpdateClientInfo request.")
	return &pb.UpdateClientInfoResponse{Success: true}, nil
}
