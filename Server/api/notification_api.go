package api

import (
	context "context"
	"errors"
	"time"

	"api.mooody.me/common"
	"api.mooody.me/models"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *MoodyAPIServer) BroadcastNotification(event *models.Notification) {
	s.notificationBroadcaster.Broadcast(event)
}

func (s *MoodyAPIServer) SendNotification(ctx context.Context, request *models.SendNotificationRequest) (*emptypb.Empty, error) {
	client, err := common.GetClientFromAuth(ctx, request.Auth, false)
	if err != nil {
		common.LogClientWithError(client, err)
		return &emptypb.Empty{}, err
	}

	if request.Notification == nil {
		return nil, errors.New("invalid request")
	}

	common.LogClient(client, `sends "[%s]: %s"`, request.Notification.Title, request.Notification.Content)

	s.BroadcastNotification(request.Notification)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeNotifications(request *models.SubscribeNotificationsRequest, server models.MoodyAPIService_SubscribeNotificationsServer) error {
	client, err := common.GetClientFromAuth(context.Background(), request.Auth, false)
	if err != nil {
		common.LogClientWithError(client, err)
		return err
	}

	common.LogClient(client, `subscribed notifications`)

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
				resp := signal.(*models.Notification)
				server.Send(resp)
			}
		case <-server.Context().Done():
			{
				common.LogClient(client, `disconnected`)
				break done
			}
		}
	}

	s.notificationBroadcaster.Unsubscribe(subscribeId)

	return nil
}
