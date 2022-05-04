package api

import (
	context "context"
	"errors"
	"time"

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

func (s *MoodyAPIServer) SubscribeNotifications(request *notifications.SubscribeRequest, server models.MoodyAPIService_SubscribeNotificationsServer) error {
	client, err := db.GetClientFromAuth(context.Background(), request.Auth, false)
	if err != nil {
		common.LogClientError(context.Background(), client, err)
		return err
	}

	common.LogClientOperation(context.Background(), client, `subscribed notifications`)

	subscribeId := time.Now().UnixNano()
	eventChannel, err := s.notificationBroadcaster.Subscribe(subscribeId)
	if err != nil {
		return err
	}

done:
	for {
		select {
		case signal := <-eventChannel:
			{
				n := signal.(*notifications.Notification)
				response := &notifications.SubscribeResponse{Notification: n}
				server.Send(response)
			}
		case <-server.Context().Done():
			{
				common.LogClientOperation(context.Background(), client, `disconnected`)
				break done
			}
		}
	}

	s.notificationBroadcaster.Unsubscribe(subscribeId)

	return nil
}
