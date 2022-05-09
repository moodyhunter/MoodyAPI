package api

import (
	"context"
	"errors"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MoodyAPIServer) BroadcastCameraEvent(event *models.CameraState) {
	s.lastCameraState = event
	s.cameraEventBroadcaster.Broadcast(event)
}

func (s *MoodyAPIServer) UpdateCameraState(ctx context.Context, request *models.UpdateCameraStateRequest) (*emptypb.Empty, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, err
	}

	if request.State == nil {
		common.LogClientOperation(ctx, client, "sent an invalid state")
		return nil, errors.New("invalid request")
	}

	common.LogClientOperation(ctx, client, "set camera state to %t", request.State.State)

	s.BroadcastCameraEvent(request.State)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeCameraStateChange(request *models.SubscribeCameraStateChangeRequest, server models.MoodyAPIService_SubscribeCameraStateChangeServer) error {
	client, err := db.GetClientFromAuth(context.Background(), request.Auth, false)
	if err != nil {
		common.LogClientError(context.Background(), client, err)
		return err
	}

	common.LogClientOperation(context.Background(), client, "subscribed to camera change event")

	server.Send(s.lastCameraState)
	s.cameraEventBroadcaster.BlockedSubscribeWithCallback(func(signal interface{}) {
		resp := signal.(*models.CameraState)
		server.Send(resp)
	})

	common.LogClientOperation(context.Background(), client, "disconnected")

	return nil
}
