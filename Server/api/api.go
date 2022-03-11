package api

import (
	context "context"
	"errors"
	"log"
	"net"
	"time"

	"api.mooody.me/broadcaster"
	"api.mooody.me/common"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type APIServer struct {
	UnimplementedMoodyAPIServiceServer
	cameraEventBroadcaster  *broadcaster.Broadcaster
	lastCameraState         *CameraState
	notificationBroadcaster *broadcaster.Broadcaster
}

func CreateAPIServer() *APIServer {
	server := &APIServer{}
	server.lastCameraState = new(CameraState)
	server.cameraEventBroadcaster = broadcaster.NewBroadcaster()
	server.notificationBroadcaster = broadcaster.NewBroadcaster()
	return server
}

func (s *APIServer) BroadcastCameraEvent(event *CameraState) {
	s.lastCameraState = event
	s.cameraEventBroadcaster.Broadcast(event)
}

func (s *APIServer) BroadcastNotification(event *Notification) {
	s.notificationBroadcaster.Broadcast(event)
}

func (s *APIServer) SubscribeCameraStateChange(request *SubscribeCameraStateChangeRequest, server MoodyAPIService_SubscribeCameraStateChangeServer) error {
	log.Printf("New gRPC camera API client connected")
	if request == nil || request.Auth == nil || request.Auth.Secret != common.APISecret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.Secret)
		return errors.New("error: Invalid Secret")
	}

	server.Send(s.lastCameraState)

	subscribeId := time.Now().UnixNano()
	eventChannel, err := s.cameraEventBroadcaster.Subscribe(subscribeId)
	if err != nil {
		return err
	}

done:
	for {
		select {
		case signal := <-eventChannel:
			{
				resp := signal.(*CameraState)
				server.Send(resp)
			}
		case <-server.Context().Done():
			{
				log.Printf("Client disconnected")
				break done
			}
		}
	}

	s.cameraEventBroadcaster.Unsubscribe(subscribeId)

	return nil
}

func (s *APIServer) SubscribeNotifications(request *SubscribeNotificationsRequest, server MoodyAPIService_SubscribeNotificationsServer) error {
	log.Printf("New notification client connected")
	if request == nil || request.Auth == nil || request.Auth.Secret != common.APISecret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.Secret)
		return errors.New("error: Invalid Secret")
	}

	subscribeId := time.Now().UnixNano()
	eventChannel, err := s.notificationBroadcaster.Subscribe(subscribeId)
	if err != nil {
		return err
	}

done:
	for {
		select {
		case signal := <-eventChannel:
			{
				resp := signal.(*Notification)
				server.Send(resp)
			}
		case <-server.Context().Done():
			{
				log.Printf("Client disconnected")
				break done
			}
		}
	}

	s.notificationBroadcaster.Unsubscribe(subscribeId)

	return nil
}

func (s *APIServer) SetCameraState(ctx context.Context, request *SetCameraStateRequest) (*emptypb.Empty, error) {
	if request == nil || request.Auth == nil || request.Auth.Secret != common.APISecret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.Secret)
		return nil, errors.New("error: Invalid Secret")
	}

	log.Printf("Changing camera state to: %t", *request.State.NewState)

	s.BroadcastCameraEvent(request.State)
	return new(emptypb.Empty), nil
}

func (s *APIServer) SendNotification(_ context.Context, request *SendNotificationRequest) (*emptypb.Empty, error) {
	if request == nil || request.Auth == nil || request.Auth.Secret != common.APISecret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.Secret)
		return &emptypb.Empty{}, errors.New("error: Invalid Secret")
	}

	println("Send Notification: ", request.Notification.Title, request.Notification.Message)
	s.BroadcastNotification(request.Notification)
	return &emptypb.Empty{}, nil
}

func StartAPIServer(listenAddress string) *APIServer {
	log.Printf("gRPC server started on: %s", listenAddress)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("WRONG")
	}

	apiServer := CreateAPIServer()
	grpcServer := grpc.NewServer()
	RegisterMoodyAPIServiceServer(grpcServer, apiServer)
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
