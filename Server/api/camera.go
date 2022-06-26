package api

import (
	"context"
	"errors"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Broadcast camera state to controllers
func (s *MoodyAPIServer) BroadcastCameraStateToControllers(event *models.CameraState) {
	s.lastCameraState = event
	s.cameraStateReportStream.Broadcast(event)
}

// Broadcast camera control signal to agents
func (s *MoodyAPIServer) BroadcastCameraControlSignalToAgents(event *models.CameraState) {
	s.lastCameraControlSignal = event
	s.cameraControlSignalStream.Broadcast(event)
}

func (s *MoodyAPIServer) SetCameraState(ctx context.Context, request *models.UpdateCameraStateRequest) (*emptypb.Empty, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, err
	}

	if request.State == nil {
		common.LogClientOperation(ctx, client, "sent an invalid state")
		return nil, errors.New("invalid request")
	}

	common.LogClientOperation(ctx, client, "set camera state to %t", request.State.State)
	s.BroadcastCameraControlSignalToAgents(request.State)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeCameraStateReport(request *models.SubscribeCameraStateChangeRequest, server models.MoodyAPIService_SubscribeCameraStateReportServer) error {
	client, err := db.AuthenticateClient(context.Background(), request.Auth, false)
	if err != nil {
		common.LogClientError(context.Background(), client, err)
		return err
	}

	common.LogClientOperation(context.Background(), client, "subscribed to camera change event")

	server.Send(s.lastCameraState)
	s.cameraStateReportStream.BlockedSubscribeWithCallback(func(signal *models.CameraState) {
		server.Send(signal)
	})

	common.LogClientOperation(context.Background(), client, "disconnected")
	return nil
}

func (s *MoodyAPIServer) ReportCameraState(ctx context.Context, request *models.UpdateCameraStateRequest) (*emptypb.Empty, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, err
	}

	if request.State == nil {
		common.LogClientOperation(ctx, client, "sent an invalid state")
		return nil, errors.New("invalid request")
	}

	common.LogClientOperation(ctx, client, "reported camera state as %t", request.State.State)
	s.BroadcastCameraStateToControllers(request.State)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeCameraControlSignal(req *models.SubscribeCameraStateChangeRequest, server models.MoodyAPIService_SubscribeCameraControlSignalServer) error {
	client, err := db.AuthenticateClient(context.Background(), req.Auth, false)
	if err != nil {
		common.LogClientError(context.Background(), client, err)
		return err
	}

	common.LogClientOperation(context.Background(), client, "subscribed to camera control signal")

	server.Send(s.lastCameraControlSignal)
	s.cameraControlSignalStream.BlockedSubscribeWithCallback(func(signal *models.CameraState) {
		server.Send(signal)
	})

	common.LogClientOperation(context.Background(), client, "disconnected")
	return nil
}
