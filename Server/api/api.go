package api

import (
	"log"
	"net"
	"time"

	"api.mooody.me/broadcaster"
	"api.mooody.me/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var moodyAPIServer *MoodyAPIServer

type MoodyAPIServer struct {
	models.UnsafeMoodyAPIServiceServer
	cameraEventStream  *broadcaster.Broadcaster
	notificationStream *broadcaster.Broadcaster
	lastCameraState    *models.CameraState

	gRPCServer    *grpc.Server
	listenAddress string
}

func CreateServer(listenAddress string) *MoodyAPIServer {
	apiServer := &MoodyAPIServer{}
	apiServer.lastCameraState = new(models.CameraState)
	apiServer.cameraEventStream = broadcaster.NewBroadcaster()
	apiServer.notificationStream = broadcaster.NewBroadcaster()
	apiServer.listenAddress = listenAddress
	log.Printf("Creating API Server on %s", listenAddress)

	apiServer.gRPCServer = grpc.NewServer()
	models.RegisterMoodyAPIServiceServer(apiServer.gRPCServer, apiServer)

	// Register reflection service on gRPC server.
	reflection.Register(apiServer.gRPCServer)
	return apiServer
}

func (apiServer *MoodyAPIServer) Serve() {
	listener, err := net.Listen("tcp", apiServer.listenAddress)
	if err != nil {
		log.Fatalf("Failed to start API Server, %s", err)
	}
	go func() {
		for {
			time.Sleep(30 * time.Second)
			apiServer.BroadcastCameraEvent(apiServer.lastCameraState)
		}
	}()
	log.Printf("API Server started on %s", apiServer.listenAddress)
	apiServer.gRPCServer.Serve(listener)
}

func (apiServer *MoodyAPIServer) Stop() {
	apiServer.gRPCServer.Stop()
	log.Println("API Server stopped")
}
