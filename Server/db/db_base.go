package db

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"api.mooody.me/models"
	"api.mooody.me/models/common"
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

func SetupConnection(dbAddress string, dbName string, dbUser string, dbPass string) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(dbAddress),
		pgdriver.WithDatabase(dbName),
		pgdriver.WithUser(dbUser),
		pgdriver.WithPassword(dbPass),
		pgdriver.WithApplicationName("MoodyAPI Server"),
		pgdriver.WithInsecure(true),
	)
	log.Printf("Connecting to database %s@%s", dbName, dbAddress)
	database = bun.NewDB(sql.OpenDB(pgconn), pgdialect.New())
	log.Println("Database connection established.")
}

func ShutdownConnection() {
	database.Close()
	log.Println("Database connection closed.")
}

func GetClientFromAuth(ctx context.Context, auth *common.Auth, requirePrivileged bool) (*common.APIClient, error) {
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

func UpdateClientLastSeen(ctx context.Context, client *common.APIClient) (*common.APIClient, error) {
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
