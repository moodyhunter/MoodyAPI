package api

import (
	"log"
	"net"
	"time"

	"api.mooody.me/broadcaster"
	"api.mooody.me/models"
	"api.mooody.me/models/light"
	"api.mooody.me/models/notifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var APIServer *MoodyAPIServer

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

	// to be used by the light APIs
	lastLightState     *light.LightState                          // the last state of the light reported by the agent
	lightControlStream *broadcaster.Broadcaster[light.LightState] // stream of light control signals [controllers => agent]

	gRPCServer    *grpc.Server
	listenAddress string
}

func CreateServer(listenAddress string) *MoodyAPIServer {
	APIServer = &MoodyAPIServer{}

	APIServer.lastCameraState = new(models.CameraState)
	APIServer.lastCameraControlSignal = new(models.CameraState)
	APIServer.lastLightState = new(light.LightState)
	APIServer.lastLightState.Brightness = 255
	APIServer.cameraStateReportStream = broadcaster.NewBroadcaster[models.CameraState]()
	APIServer.cameraControlSignalStream = broadcaster.NewBroadcaster[models.CameraState]()
	APIServer.notificationStream = broadcaster.NewBroadcaster[notifications.Notification]()
	APIServer.keepAliveStream = broadcaster.NewBroadcaster[models.KeepAliveMessage]()
	APIServer.lightControlStream = broadcaster.NewBroadcaster[light.LightState]()

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

	go func() {
		for {
			time.Sleep(30 * time.Second)
			apiServer.keepAliveStream.Broadcast(&models.KeepAliveMessage{Time: timestamppb.Now()})
			apiServer.cameraControlSignalStream.Broadcast(apiServer.lastCameraControlSignal)
			apiServer.cameraStateReportStream.Broadcast(apiServer.lastCameraState)
		}
	}()

	log.Printf("API Server started on %s", apiServer.listenAddress)
	apiServer.gRPCServer.Serve(listener)
}

func (apiServer *MoodyAPIServer) Stop() {
	apiServer.gRPCServer.Stop()
	log.Println("API Server stopped")
}
