package api

import (
	"time"

	"api.mooody.me/broadcaster"
	"api.mooody.me/models"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MoodyAPIServer struct {
	models.UnimplementedMoodyAPIServiceServer
	cameraEventBroadcaster  *broadcaster.Broadcaster
	notificationBroadcaster *broadcaster.Broadcaster
	lastCameraState         *models.CameraState
}

func CreateServer() (*MoodyAPIServer, *grpc.Server) {
	apiServer := &MoodyAPIServer{}
	apiServer.lastCameraState = new(models.CameraState)
	apiServer.cameraEventBroadcaster = broadcaster.NewBroadcaster()
	apiServer.notificationBroadcaster = broadcaster.NewBroadcaster()

	grpcServer := grpc.NewServer()
	models.RegisterMoodyAPIServiceServer(grpcServer, apiServer)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	go func() {
		for {
			time.Sleep(30 * time.Second)
			apiServer.BroadcastCameraEvent(apiServer.lastCameraState)
		}
	}()

	return apiServer, grpcServer
}
