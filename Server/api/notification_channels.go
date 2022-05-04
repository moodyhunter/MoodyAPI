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

	// Log current client's operation
	common.LogClientOperation(ctx, client, `creates notification channel`)

	return nil, nil
}

func (s *MoodyAPIServer) DeleteNotificationChannel(ctx context.Context, request *notifications.DeleteChannelRequest) (*notifications.DeleteChannelResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	// log current client's operation
	common.LogClientOperation(ctx, client, `deletes notification channel`)

	return nil, nil
}

func (s *MoodyAPIServer) ListNotificationChannel(ctx context.Context, request *notifications.ListChannelRequest) (*notifications.ListChannelResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	// log current client's operation
	common.LogClientOperation(ctx, client, `lists notification channels`)

	return nil, nil
}

func (s *MoodyAPIServer) UpdateNotificationChannel(ctx context.Context, request *notifications.UpdateChannelRequest) (*notifications.UpdateChannelResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	// log current client's operation
	common.LogClientOperation(ctx, client, `updates notification channel`)

	return nil, nil
}
