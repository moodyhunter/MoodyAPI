package db

import (
	"context"
	"errors"

	"api.mooody.me/models/notifications"
)

func GetNotificationChannelById(ctx context.Context, channelId int64) (*notifications.NotificationChannel, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	channelORM := &notifications.NotificationChannelORM{}

	err = database.NewSelect().
		Model(channelORM).
		Where("id = ?", channelId).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return channelORM.ToPB(ctx)
}

func ListNotificationChannels(ctx context.Context) ([]*notifications.NotificationChannel, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	channels := []*notifications.NotificationChannelORM{}

	err = database.NewSelect().
		Model(&channels).
		Order("id").
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	result := []*notifications.NotificationChannel{}
	for _, v := range channels {
		pbObject, pbErr := v.ToPB(ctx)
		if pbErr != nil {
			return nil, err
		}
		result = append(result, pbObject)
	}

	return result, nil
}

func CreateNotificationChannel(ctx context.Context, channel *notifications.NotificationChannel) (*notifications.NotificationChannel, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	channelORM, err := channel.ToORM(ctx)
	if err != nil {
		return nil, err
	}

	result, err := database.NewInsert().
		Model(&channelORM).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	r, err := result.RowsAffected()
	if err != nil {
		return nil, err
	} else if r == 0 {
		return nil, errors.New("unexpected affected rows")
	}

	return channelORM.ToPB(ctx)
}

func UpdateNotificationChannel(ctx context.Context, channel *notifications.NotificationChannel) (*notifications.NotificationChannel, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	channelORM, err := channel.ToORM(ctx)
	if err != nil {
		return nil, err
	}

	result, err := database.NewUpdate().
		Model(&channelORM).
		WherePK().
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	r, err := result.RowsAffected()
	if err != nil {
		return nil, err
	} else if r == 0 {
		return nil, errors.New("unexpected affected rows")
	}

	return channelORM.ToPB(ctx)
}

func DeleteNotificationChannel(ctx context.Context, channelId int64) (*notifications.NotificationChannel, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	channelORM := &notifications.NotificationChannelORM{}

	result, err := database.NewDelete().
		Model(&channelORM).
		Where("id = ?", channelId).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	r, err := result.RowsAffected()
	if err != nil {
		return nil, err
	} else if r == 0 {
		return nil, errors.New("unexpected affected rows")
	}

	return channelORM.ToPB(ctx)

}
