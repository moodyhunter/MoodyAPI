package camapi

import (
	context "context"
	"errors"
	"log"
	"net"
	"time"

	"api.mooody.me/broadcaster"
	"api.mooody.me/common"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type CamAPIServer struct {
	UnimplementedCameraServiceServer
	eventSender *broadcaster.Broadcaster
	lastEvent   *CameraState
}

func NewCamAPIServer() *CamAPIServer {
	server := &CamAPIServer{}
	server.lastEvent = new(CameraState)
	server.eventSender = broadcaster.NewBroadcaster()
	return server
}

func (s *CamAPIServer) BroadcaseEvent(event *CameraState) {
	s.lastEvent = event
	s.eventSender.Broadcast(event)
}

func (s *CamAPIServer) SubscribeCameraStateChange(request *SubscribeCameraStateChangeRequest, server CameraService_SubscribeCameraStateChangeServer) error {
	log.Printf("New gRPC camera API client connected")
	if request == nil || request.Auth == nil || request.Auth.Secret != common.APISecret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.Secret)
		return errors.New("error: Invalid Secret")
	}

	server.Send(s.lastEvent)

	subscribeId := time.Now().UnixNano()
	eventChannel, err := s.eventSender.Subscribe(subscribeId)
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

	s.eventSender.Unsubscribe(subscribeId)

	return nil
}

func (s *CamAPIServer) SetCameraState(ctx context.Context, request *SetCameraStateRequest) (*emptypb.Empty, error) {
	if request == nil || request.Auth == nil || request.Auth.Secret != common.APISecret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.Secret)
		return nil, errors.New("error: Invalid Secret")
	}

	log.Printf("Applying new camera state: %t", *request.State.NewState)

	s.BroadcaseEvent(request.State)
	return new(emptypb.Empty), nil
}

func StartAPIServer(listenAddress string) *CamAPIServer {
	log.Printf("gRPC server started on: %s", listenAddress)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("WRONG")
	}

	camAPIServer := NewCamAPIServer()
	server := grpc.NewServer()
	RegisterCameraServiceServer(server, camAPIServer)
	go server.Serve(lis)
	return camAPIServer
}
