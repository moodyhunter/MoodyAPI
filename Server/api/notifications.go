package api

import (
	context "context"
	"errors"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models"
	"api.mooody.me/models/notifications"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *MoodyAPIServer) BroadcastNotification(event *notifications.Notification) {
	s.notificationBroadcaster.Broadcast(event)
}

func (s *MoodyAPIServer) SendNotification(ctx context.Context, request *notifications.SendRequest) (*emptypb.Empty, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return &emptypb.Empty{}, err
	}

	if request.Notification == nil {
		common.LogClientOperation(ctx, client, "sent an invalid notification")
		return nil, errors.New("invalid request")
	}

	common.LogClientOperation(ctx, client, `sends "[%s]: %s"`, request.Notification.Title, request.Notification.Content)

	s.BroadcastNotification(request.Notification)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeNotificationInternal(callback func(signal *notifications.Notification)) error {
	s.notificationBroadcaster.SubscribeWithCallback(func(signal interface{}) {
		resp := signal.(*notifications.Notification)
		callback(resp)
	})
	return nil
}

func (s *MoodyAPIServer) SubscribeNotifications(request *notifications.SubscribeRequest, server models.MoodyAPIService_SubscribeNotificationsServer) error {
	client, err := db.GetClientFromAuth(context.Background(), request.Auth, false)
	if err != nil {
		common.LogClientError(context.Background(), client, err)
		return err
	}

	common.LogClientOperation(context.Background(), client, `subscribed notifications`)

	s.notificationBroadcaster.SubscribeWithCallback(func(signal interface{}) {
		resp := signal.(*notifications.Notification)
		server.Send(resp)
	})

	common.LogClientOperation(context.Background(), client, `disconnected`)

	return nil
}

func (s *MoodyAPIServer) ListNotifications(ctx context.Context, request *notifications.ListRequest) (*notifications.ListResponse, error) {
	client, err := db.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		common.LogClientError(ctx, client, err)
		return nil, err
	}

	common.LogClientOperation(ctx, client, `lists notifications`)
	result, err := db.ListNotifications(ctx, request.ChannelID, request.SenderID, request.Urgency, request.Private)
	if err != nil {
		return nil, err
	}

	return &notifications.ListResponse{Notifications: result}, nil
}
