package api

import (
	"context"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models/notifications"
)

func (s *MoodyAPIServer) CreateNotificationChannel(ctx context.Context, request *notifications.CreateChannelRequest) (*notifications.CreateChannelResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `creates notification channel`)
	result, err := db.CreateNotificationChannel(ctx, request.Channel)
	if err != nil {
		return nil, err
	}

	return &notifications.CreateChannelResponse{Channel: result}, nil
}

func (s *MoodyAPIServer) DeleteNotificationChannel(ctx context.Context, request *notifications.DeleteChannelRequest) (*notifications.DeleteChannelResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `deletes notification channel`)
	db.DeleteNotificationChannel(ctx, request.ChannelID)
	return &notifications.DeleteChannelResponse{}, nil
}

func (s *MoodyAPIServer) ListNotificationChannel(ctx context.Context, request *notifications.ListChannelRequest) (*notifications.ListChannelResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `lists notification channels`)
	result, err := db.ListNotificationChannels(ctx)
	if err != nil {
		return nil, err
	}

	return &notifications.ListChannelResponse{Channels: result}, nil
}

func (s *MoodyAPIServer) UpdateNotificationChannel(ctx context.Context, request *notifications.UpdateChannelRequest) (*notifications.UpdateChannelResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `updates notification channel`)
	result, err := db.UpdateNotificationChannel(ctx, request.Channel)
	if err != nil {
		return nil, err
	}

	return &notifications.UpdateChannelResponse{Channel: result}, nil
}
