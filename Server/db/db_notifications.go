package db

import (
	"context"

	"api.mooody.me/models/common"
	"api.mooody.me/models/notifications"
)

func ListNotifications(ctx context.Context, ChannelID *int64, SenderID *int64, Urgency *common.NotificationUrgency, Private *bool) ([]*notifications.Notification, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, channel_id, sender_id, urgency, private, title, content FROM notifications WHERE 1=1`

	if ChannelID != nil {
		query += ` AND channel_id = ?`
	}

	if SenderID != nil {
		query += ` AND sender_id = ?`
	}

	if Urgency != nil {
		query += ` AND urgency = ?`
	}

	if Private != nil {
		query += ` AND private = ?`
	}

	query += ` ORDER BY id`

	rows, err := database.QueryContext(ctx, query, ChannelID, SenderID, Urgency, Private)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*notifications.Notification
	for rows.Next() {
		var notification notifications.Notification
		err := rows.Scan(&notification.Id, &notification.ChannelId, &notification.SenderId, &notification.Urgency, &notification.Private, &notification.Title, &notification.Content)
		if err != nil {
			return nil, err
		}
		result = append(result, &notification)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}
