package api

import (
	context "context"
	"errors"
	"log"
	"time"

	"api.mooody.me/db"
	"api.mooody.me/models"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *MoodyAPIServer) BroadcastNotification(event *models.Notification) {
	s.notificationBroadcaster.Broadcast(event)
}

func (s *MoodyAPIServer) SendNotification(ctx context.Context, request *models.SendNotificationRequest) (*emptypb.Empty, error) {
	if request == nil || request.Auth == nil {
		log.Printf("bad request")
		return nil, errors.New("invalid client id")
	}

	client, err := db.GetClientByUUID(ctx, request.Auth.ClientUuid)
	if err != nil {
		log.Printf("invalid client id: %s", request.Auth.ClientUuid)
		return &emptypb.Empty{}, errors.New("invalid client id")
	}

	if !client.GetEnabled() {
		log.Printf("[%s] client is not enabled", *client.Name)
		return &emptypb.Empty{}, errors.New("client is not enabled")
	}

	log.Printf("[%s] sends notification: [%s]: %s", *client.Name, request.Notification.Title, request.Notification.Message)

	s.BroadcastNotification(request.Notification)
	return &emptypb.Empty{}, nil
}

func (s *MoodyAPIServer) SubscribeNotifications(request *models.SubscribeNotificationsRequest, server models.MoodyAPIService_SubscribeNotificationsServer) error {
	if request == nil || request.Auth == nil {
		log.Printf("bad request")
		return errors.New("invalid client id")
	}

	client, err := db.GetClientByUUID(context.Background(), request.Auth.ClientUuid)
	if err != nil {
		log.Printf("invalid client id: %s", request.Auth.ClientUuid)
		return errors.New("invalid client id")
	}

	if !client.GetEnabled() {
		log.Printf("[%s] client is not enabled", *client.Name)
		return errors.New("client is not enabled")
	}

	log.Printf("[%s] subscribes to camera change info", *client.Name)

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
				log.Printf("client %s disconnected", *client.Name)
				break done
			}
		}
	}

	s.notificationBroadcaster.Unsubscribe(subscribeId)

	return nil
}
