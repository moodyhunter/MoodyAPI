package api

import (
	"log"
	"net"
	"time"

	"api.mooody.me/api/pb"
	"api.mooody.me/broadcaster"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MoodyAPIServer struct {
	pb.UnsafeMoodyAPIServiceServer
	cameraEventBroadcaster  *broadcaster.Broadcaster
	notificationBroadcaster *broadcaster.Broadcaster
	lastCameraState         *pb.CameraState
}

func CreateAPIServer() *MoodyAPIServer {
	server := &MoodyAPIServer{}
	server.lastCameraState = new(pb.CameraState)
	server.cameraEventBroadcaster = broadcaster.NewBroadcaster()
	server.notificationBroadcaster = broadcaster.NewBroadcaster()
	return server
}

func StartAPIServer(listenAddress string) *MoodyAPIServer {
	log.Printf("gRPC server started on: %s", listenAddress)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("WRONG")
	}

	apiServer := CreateAPIServer()
	grpcServer := grpc.NewServer()
	pb.RegisterMoodyAPIServiceServer(grpcServer, apiServer)
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	go func() {
		for {
			time.Sleep(30 * time.Second)
			apiServer.BroadcastCameraEvent(apiServer.lastCameraState)
		}
	}()
	go grpcServer.Serve(lis)
	return apiServer
}
