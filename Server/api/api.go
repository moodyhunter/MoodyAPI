package api

import (
	"log"
	"net"
	"time"

	"api.mooody.me/broadcaster"
	"api.mooody.me/models"
	"api.mooody.me/models/notifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var APIServer *MoodyAPIServer

type MoodyAPIServer struct {
	models.UnsafeMoodyAPIServiceServer

	notificationStream *broadcaster.Broadcaster[notifications.Notification]

	gRPCServer    *grpc.Server
	listenAddress string
}

func CreateServer(listenAddress string) *MoodyAPIServer {
	APIServer = &MoodyAPIServer{}

	APIServer.notificationStream = broadcaster.NewBroadcaster[notifications.Notification]()

	APIServer.listenAddress = listenAddress
	log.Printf("Creating API Server on %s", listenAddress)

	TIMEOUT := time.Hour * 1
	APIServer.gRPCServer = grpc.NewServer(
		grpc.ConnectionTimeout(TIMEOUT),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: TIMEOUT,
		}),
	)
	models.RegisterMoodyAPIServiceServer(APIServer.gRPCServer, APIServer)

	// Register reflection service on gRPC server.
	reflection.Register(APIServer.gRPCServer)
	return APIServer
}

func (apiServer *MoodyAPIServer) Serve() {
	listener, err := net.Listen("tcp", apiServer.listenAddress)
	if err != nil {
		log.Fatalf("Failed to start API Server, %s", err)
	}

	log.Printf("API Server started on %s", apiServer.listenAddress)
	apiServer.gRPCServer.Serve(listener)
}

func (apiServer *MoodyAPIServer) Stop() {
	apiServer.gRPCServer.Stop()
	log.Println("API Server stopped")
}
