package db

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"api.mooody.me/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	database *bun.DB = nil
)

func checkDatabaseConnectivity() error {
	if database == nil {
		log.Fatal("Database not connected.")
		return errors.New("database not connected")
	}
	return nil
}

func GetClientFromAuth(ctx context.Context, auth *models.Auth, requirePrivileged bool) (*models.APIClient, error) {
	if auth == nil {
		return nil, errors.New("invalid Auth")
	}

	client, err := GetClientByUUID(ctx, auth.ClientUuid)
	if err != nil {
		return nil, errors.New("invalid client id")
	}

	client, err = UpdateClientLastSeen(ctx, client)
	if err != nil {
		return client, errors.New("failed to update client last seen")
	}

	if !client.GetEnabled() {
		return client, errors.New("client is not enabled")
	}

	if requirePrivileged {
		if !client.GetPrivileged() {
			return client, errors.New("client isn't privileged as required")
		}
	}

	return client, nil
}

func SetupDBConnection(dbAddress string, dbName string, dbUser string, dbPass string) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(dbAddress),
		pgdriver.WithDatabase(dbName),
		pgdriver.WithUser(dbUser),
		pgdriver.WithPassword(dbPass),
		pgdriver.WithApplicationName("MoodyAPI Server"),
		pgdriver.WithInsecure(true),
	)
	database = bun.NewDB(sql.OpenDB(pgconn), pgdialect.New())
}

func GetClientByUUID(ctx context.Context, clientUuid string) (*models.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	client := models.APIClientORM{}

	q := database.NewSelect().
		Model(&client).
		Where("client_uuid = ?", clientUuid)

	err = q.Limit(1).Scan(ctx)

	if err != nil {
		return nil, err
	}

	return client.ToPB(ctx)
}

func GetClientByID(ctx context.Context, id int64) (*models.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	clientORM := models.APIClientORM{}
	err = database.NewSelect().
		Model(&clientORM).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return clientORM.ToPB(ctx)
}

func UpdateClientLastSeen(ctx context.Context, client *models.APIClient) (*models.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return client, err
	}

	newClient, err := GetClientByID(ctx, client.Id)
	newClient.LastSeen = timestamppb.Now()

	err = UpdateClient(ctx, newClient)
	if err != nil {
		return client, err
	}
	return newClient, nil
}

func CreateClient(ctx context.Context, client *models.APIClient) (*models.APIClient, error) {
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

func UpdateClient(ctx context.Context, client *models.APIClient) error {
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

func DeleteClient(ctx context.Context, client *models.APIClient) error {
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

func ListClients(ctx context.Context) ([]*models.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	clientORM := []models.APIClientORM{}
	err = database.NewSelect().
		Model(&clientORM).
		Scan(ctx)

	clients := []*models.APIClient{}
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

func LogOperation(ctx context.Context, logInfo *models.OperationLog) error {
	logORM, _ := logInfo.ToORM(ctx)
	res, err := database.NewInsert().
		Model(&logORM).
		Exec(ctx)

	if err != nil {
		return err
	}

	r, err := res.RowsAffected()
	if err != nil {
		return err
	} else if r == 0 {
		return errors.New("unexpected affected rows")
	}

	return nil
}
