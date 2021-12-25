package camapi

import (
	context "context"
	"errors"
	"fmt"
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
}

func NewCamAPIServer() *CamAPIServer {
	server := &CamAPIServer{}
	server.eventSender = broadcaster.NewBroadcaster()
	return server
}

func (s *CamAPIServer) BroadcaseEvent(resp *CameraStateChangedResponse) {
	s.eventSender.Broadcast(resp)
}

func (s *CamAPIServer) SubscribeCameraStateChange(request *SubscribeCameraStateChangeRequest, server CameraService_SubscribeCameraStateChangeServer) error {
	log.Printf("New gRPC camera API client connected")
	if request == nil || request.Auth == nil || request.Auth.Secret != common.Secret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.Secret)
		return errors.New("error: Invalid Eecret")
	}

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
				fmt.Println("Responsing with new changed response")
				resp := signal.(*CameraStateChangedResponse)
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

func (s *CamAPIServer) SetCameraState(ctx context.Context, req *SetCameraStateRequest) (*emptypb.Empty, error) {
	return nil, nil
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

	go func() {
		for {
			time.Sleep(1 * time.Second)
			camAPIServer.BroadcaseEvent(&CameraStateChangedResponse{
				IsInitial: false,
				Values:    &CameraStateChangedResponse_MotionEventId{"test"},
			})
		}
	}()

	return camAPIServer
}
