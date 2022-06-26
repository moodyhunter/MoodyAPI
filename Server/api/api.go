package api

import (
	"log"
	"net"
	"time"

	"api.mooody.me/broadcaster"
	"api.mooody.me/models"
	"api.mooody.me/models/notifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var moodyAPIServer *MoodyAPIServer

type MoodyAPIServer struct {
	models.UnsafeMoodyAPIServiceServer

	notificationStream *broadcaster.Broadcaster[notifications.Notification]
	keepAliveStream    *broadcaster.Broadcaster[models.KeepAliveMessage]

	// to be used by the agent APIs
	lastCameraState         *models.CameraState                          // the last state of the camera reported by the agent
	cameraStateReportStream *broadcaster.Broadcaster[models.CameraState] // stream of camera state changes [agent => controllers]

	// to be used by the controllers APIs
	lastCameraControlSignal   *models.CameraState                          // the last state of the camera control signal sent by the controllers
	cameraControlSignalStream *broadcaster.Broadcaster[models.CameraState] // stream of control signals [controllers => agent]

	gRPCServer    *grpc.Server
	listenAddress string
}

func CreateServer(listenAddress string) *MoodyAPIServer {
	apiServer := &MoodyAPIServer{}

	apiServer.lastCameraState = new(models.CameraState)
	apiServer.lastCameraControlSignal = new(models.CameraState)
	apiServer.cameraStateReportStream = broadcaster.NewBroadcaster[models.CameraState]()
	apiServer.cameraControlSignalStream = broadcaster.NewBroadcaster[models.CameraState]()
	apiServer.notificationStream = broadcaster.NewBroadcaster[notifications.Notification]()
	apiServer.keepAliveStream = broadcaster.NewBroadcaster[models.KeepAliveMessage]()

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
			apiServer.keepAliveStream.Broadcast(&models.KeepAliveMessage{Time: timestamppb.Now()})
		}
	}()

	log.Printf("API Server started on %s", apiServer.listenAddress)
	apiServer.gRPCServer.Serve(listener)
}

func (apiServer *MoodyAPIServer) Stop() {
	apiServer.gRPCServer.Stop()
	log.Println("API Server stopped")
}
