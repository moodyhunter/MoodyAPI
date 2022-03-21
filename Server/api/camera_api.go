package api

import (
	context "context"
	"errors"
	"log"
	"time"

	"api.mooody.me/api/pb"
	"api.mooody.me/common"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *MoodyAPIServer) BroadcastCameraEvent(event *pb.CameraState) {
	s.lastCameraState = event
	s.cameraEventBroadcaster.Broadcast(event)
}

func (s *MoodyAPIServer) UpdateCameraState(ctx context.Context, request *pb.UpdateCameraStateRequest) (*emptypb.Empty, error) {
	if request == nil || request.Auth == nil || request.Auth.ClientId != common.APISecret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.ClientId)
		return nil, errors.New("error: Invalid Secret")
	}

	log.Printf("Changing camera state to: %t", request.State.GetState())

	s.BroadcastCameraEvent(request.State)
	return new(emptypb.Empty), nil
}

func (s *MoodyAPIServer) SubscribeCameraStateChange(request *pb.SubscribeCameraStateChangeRequest, server pb.MoodyAPIService_SubscribeCameraStateChangeServer) error {
	log.Printf("New gRPC camera API client connected")
	if request == nil || request.Auth == nil || request.Auth.ClientId != common.APISecret {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.ClientId)
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
				resp := signal.(*pb.CameraState)
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
