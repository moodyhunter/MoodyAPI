package api

import (
	context "context"
	"log"
	"net"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type cameraServiceServer struct {
	UnimplementedCameraServiceServer
}

func (s *cameraServiceServer) SubscribeCameraStateChange(_ *emptypb.Empty, server CameraService_SubscribeCameraStateChangeServer) error {
	// resp := CameraStateChangedResponses{}
	// resp.States = append(resp.States, &CameraState{})
	// server.Send(&resp)

	ctx := server.Context()
	done := <-ctx.Done()
	log.Printf("Client disconnected, %s", done)
	return nil
}

func (s *cameraServiceServer) SetCameraState(ctx context.Context, req *SetCameraStateRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func StartAPIServer(listenAddress string) {
	log.Printf("gRPC server started on: %s", listenAddress)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("WRONG")
	}
	server := grpc.NewServer()
	RegisterCameraServiceServer(server, &cameraServiceServer{})
	go server.Serve(lis)
}
