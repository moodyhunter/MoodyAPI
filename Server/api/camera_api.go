package api

import (
	context "context"
	"errors"
	"log"
	"time"

	"api.mooody.me/db"
	"api.mooody.me/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MoodyAPIServer) BroadcastCameraEvent(event *models.CameraState) {
	s.lastCameraState = event
	s.cameraEventBroadcaster.Broadcast(event)
}

func (s *MoodyAPIServer) UpdateCameraState(ctx context.Context, request *models.UpdateCameraStateRequest) (*emptypb.Empty, error) {
	if request == nil || request.Auth == nil {
		log.Printf("bad request")
		return nil, errors.New("invalid client id")
	}

	client, err := db.GetClientByUUID(ctx, request.Auth.ClientUuid)
	if err != nil {
		log.Printf("invalid client id: %s", request.Auth.ClientUuid)
		return &emptypb.Empty{}, errors.New("invalid client id")
	}

	if !client.GetEnabled() {
		log.Printf("[%s] client is not enabled", *client.Name)
		return &emptypb.Empty{}, errors.New("client is not enabled")
	}

	log.Printf("[%s] sets camera state to %t", *client.Name, request.State.State)

	s.BroadcastCameraEvent(request.State)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeCameraStateChange(request *models.SubscribeCameraStateChangeRequest, server models.MoodyAPIService_SubscribeCameraStateChangeServer) error {
	if request == nil || request.Auth == nil {
		log.Printf("bad request")
		return errors.New("invalid client id")
	}

	client, err := db.GetClientByUUID(context.Background(), request.Auth.ClientUuid)
	if err != nil {
		log.Printf("invalid client id: %s", request.Auth.ClientUuid)
		return errors.New("invalid client id")
	}

	if !client.GetEnabled() {
		log.Printf("[%s] client is not enabled", *client.Name)
		return errors.New("client is not enabled")
	}

	log.Printf("[%s] subscribes to camera change info", *client.Name)

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
				resp := signal.(*models.CameraState)
				server.Send(resp)
			}
		case <-server.Context().Done():
			{
				log.Printf("client %s disconnected", *client.Name)
				break done
			}
		}
	}

	s.cameraEventBroadcaster.Unsubscribe(subscribeId)

	return nil
}
