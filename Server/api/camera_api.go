package api

import (
	"context"
	"errors"
	"time"

	"api.mooody.me/common"
	"api.mooody.me/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MoodyAPIServer) BroadcastCameraEvent(event *models.CameraState) {
	s.lastCameraState = event
	s.cameraEventBroadcaster.Broadcast(event)
}

func (s *MoodyAPIServer) UpdateCameraState(ctx context.Context, request *models.UpdateCameraStateRequest) (*emptypb.Empty, error) {
	client, err := common.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		common.LogClientWithError(client, err)
		return nil, err
	}

	if request.State == nil {
		return nil, errors.New("invalid request")
	}

	common.LogClient(client, "set camera state to %t", request.State.State)

	s.BroadcastCameraEvent(request.State)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeCameraStateChange(request *models.SubscribeCameraStateChangeRequest, server models.MoodyAPIService_SubscribeCameraStateChangeServer) error {
	client, err := common.GetClientFromAuth(context.Background(), request.Auth, false)
	if err != nil {
		common.LogClientWithError(client, err)
		return err
	}

	common.LogClient(client, "subscribed to camera change event")

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
				common.LogClient(client, "disconnected", *client.Name)
				break done
			}
		}
	}

	s.cameraEventBroadcaster.Unsubscribe(subscribeId)

	return nil
}
