package api

import (
	"context"
	"errors"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models"
	"api.mooody.me/models/light"
)

func (s *MoodyAPIServer) BroadcastLightState(status *light.LightState) {
	// preprocess color with brightness if warmwhite is not set
	if !status.GetWarmwhite() {
		if status.GetColored() == nil {
			status.Mode = &light.LightState_Warmwhite{Warmwhite: true}
		} else {
			var red = uint32(float32(status.GetColored().Red) * (float32(status.GetBrightness()) / 255))
			var green = uint32(float32(status.GetColored().Green) * (float32(status.GetBrightness()) / 255))
			var blue = uint32(float32(status.GetColored().Blue) * (float32(status.GetBrightness()) / 255))

			status.GetColored().Blue = blue
			status.GetColored().Green = green
			status.GetColored().Red = red
		}
	}
	s.LastLightState = status
	s.lightControlStream.Broadcast(status)
}

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

	s.BroadcastLightState(request.State)
	return &light.SetLightResponse{}, nil
}

func (s *MoodyAPIServer) GetLightState(ctx context.Context, request *light.GetLightRequest) (*light.GetLightResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, err
	}

	return &light.GetLightResponse{State: s.LastLightState}, nil
}

func (s *MoodyAPIServer) SubscribeLightStateChange(subscribeLightStateRequest *light.SubscribeLightRequest, server models.MoodyAPIService_SubscribeLightStateChangeServer) error {
	client, err := db.AuthenticateClient(server.Context(), subscribeLightStateRequest.Auth, false)
	if err != nil {
		common.LogClientError(server.Context(), client, err)
		return err
	}

	server.Send(s.LastLightState)
	s.lightControlStream.BlockedSubscribeWithCallback(func(signal *light.LightState) {
		server.Send(signal)
	})

	common.LogClientOperation(context.Background(), client, "disconnected")
	return nil
}
