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
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.ClientUuid)
		return nil, errors.New("invalid client id")
	}

	client, err := db.GetClientByUUID(ctx, request.Auth.ClientUuid)
	if err != nil {
		return &emptypb.Empty{}, errors.New("invalid client id")
	}

	log.Printf("client %s (%s) sets camera state to %s", *client.Name, *client.Uuid, request.State.State)

	s.BroadcastCameraEvent(request.State)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeCameraStateChange(request *models.SubscribeCameraStateChangeRequest, server models.MoodyAPIService_SubscribeCameraStateChangeServer) error {
	if request == nil || request.Auth == nil {
		log.Printf("WARNING: Invalid secret from client: %s", request.Auth.ClientUuid)
		return errors.New("invalid client id")
	}

	client, err := db.GetClientByUUID(context.Background(), request.Auth.ClientUuid)
	if err != nil {
		return errors.New("invalid client id")
	}

	log.Printf("client %s (%s) subscribes to camera change info", *client.Name, *client.Uuid)

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
