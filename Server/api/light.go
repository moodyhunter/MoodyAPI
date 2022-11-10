package api

import (
	"context"
	"errors"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models"
	"api.mooody.me/models/light"
)

func (s *MoodyAPIServer) SetLightState(ctx context.Context, request *light.SetLightRequest) (*light.SetLightResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, err
	}

	if request.State == nil {
		err = errors.New("state is nil")
		common.LogClientError(ctx, client, err)
		return nil, err
	}

	common.LogClientOperation(ctx, client, "set light to %v", request.State.On)

	s.lightControlStream.Broadcast(request.State)
	s.lastLightState = request.State
	return &light.SetLightResponse{}, nil
}

func (s *MoodyAPIServer) GetLightState(ctx context.Context, request *light.GetLightRequest) (*light.GetLightResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, err
	}

	common.LogClientOperation(ctx, client, "get light state")

	return &light.GetLightResponse{State: s.lastLightState}, nil
}

func (s *MoodyAPIServer) SubscribeLightStateChange(subscribeLightStateRequest *light.SubscribeLightRequest, server models.MoodyAPIService_SubscribeLightStateChangeServer) error {
	client, err := db.AuthenticateClient(server.Context(), subscribeLightStateRequest.Auth, false)
	if err != nil {
		common.LogClientError(server.Context(), client, err)
		return err
	}

	server.Send(s.lastLightState)
	s.lightControlStream.BlockedSubscribeWithCallback(func(signal *light.LightState) {
		server.Send(signal)
	})

	common.LogClientOperation(context.Background(), client, "disconnected")
	return nil
}
