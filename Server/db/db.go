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

func CheckClientValidity(context context.Context, clientUuid string, requirePrivileged bool) (*models.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	client := models.APIClientORM{}

	q := database.NewSelect().
		Model(&client).
		Where("client_uuid = ?", clientUuid)

	if requirePrivileged {
		q = q.Where("privileged = true")
	}

	err = q.Limit(1).Scan(context)

	if err != nil {
		return nil, err
	}

	return client.ToPB(context)
}

func ListClients(context context.Context) ([]*models.APIClient, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	clientORM := []models.APIClientORM{}
	err = database.NewSelect().Model(&clientORM).Scan(context)

	clients := []*models.APIClient{}
	if err != nil {
		return nil, err
	}

	for _, v := range clientORM {
		pbObject, pbErr := v.ToPB(context)
		if pbErr != nil {
			return nil, err
		}
		clients = append(clients, pbObject)
	}

	return clients, nil

}
