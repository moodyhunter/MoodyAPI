package db

import (
	"context"
	"errors"

	"api.mooody.me/models/common"
)

func GetClientByUUID(ctx context.Context, clientUuid string) (*common.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	client := common.APIClientORM{}

	q := database.NewSelect().
		Model(&client).
		Where("client_uuid = ?", clientUuid)

	err = q.Limit(1).Scan(ctx)

	if err != nil {
		return nil, err
	}

	return client.ToPB(ctx)
}

func GetClientByID(ctx context.Context, id int64) (*common.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	clientORM := common.APIClientORM{}
	err = database.NewSelect().
		Model(&clientORM).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return clientORM.ToPB(ctx)
}

func CreateClient(ctx context.Context, client *common.APIClient) (*common.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	clientORM, err := client.ToORM(ctx)
	if err != nil {
		return nil, err
	}

	database.ExecContext(ctx, "")

	result, err := database.NewInsert().
		Model(&clientORM).
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

	return clientORM.ToPB(ctx)
}

func ListClients(ctx context.Context) ([]*common.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	clientORM := []common.APIClientORM{}
	err = database.NewSelect().
		Model(&clientORM).
		Scan(ctx)

	clients := []*common.APIClient{}
	if err != nil {
		return nil, err
	}

	for _, v := range clientORM {
		pbObject, pbErr := v.ToPB(ctx)
		if pbErr != nil {
			return nil, err
		}
		clients = append(clients, pbObject)
	}

	return clients, nil
}

func UpdateClient(ctx context.Context, client *common.APIClient) error {
	err := checkDatabaseConnectivity()
	if err != nil {
		return err
	}

	clientORM, err := client.ToORM(ctx)

	if err != nil {
		return err
	}

	result, err := database.NewUpdate().
		Model(&clientORM).
		WherePK().
		Exec(ctx)

	if err != nil {
		return err
	}

	r, err := result.RowsAffected()
	if err != nil {
		return err
	} else if r == 0 {
		return errors.New("unexpected affected rows")
	}

	return nil
}

func DeleteClient(ctx context.Context, client *common.APIClient) error {
	err := checkDatabaseConnectivity()
	if err != nil {
		return err
	}

	clientORM, err := client.ToORM(ctx)

	if err != nil {
		return err
	}

	result, err := database.NewDelete().
		Model(&clientORM).
		WherePK().
		Exec(ctx)

	r, err := result.RowsAffected()
	if err != nil {
		return err
	} else if r == 0 {
		return errors.New("unexpected affected rows")
	}
	return nil
}
