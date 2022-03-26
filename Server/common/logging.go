package common

import (
	"context"
	"fmt"
	"log"

	"api.mooody.me/db"
	"api.mooody.me/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getClientIdName(client *models.APIClient) (int64, string) {
	var cid int64
	var cname string

	if client != nil {
		cid = client.Id
		cname = *client.Name
	} else {
		cid = 0
		cname = "<unknown>"
	}
	return cid, cname
}

func writeLog(ctx context.Context, client *models.APIClient, operation string) {
	cid, cname := getClientIdName(client)
	log.Printf(operation)

	logInfo := models.OperationLog{
		ClientId:   cid,
		ClientName: cname,
		Time:       timestamppb.Now(),
		Operation:  operation,
	}

	db.LogOperation(ctx, &logInfo)
}

func LogClientOperation(ctx context.Context, client *models.APIClient, message string, args ...interface{}) {
	frame := GetCallerFunctionName()
	writeLog(ctx, client, "I: "+frame+": "+fmt.Sprintf(message, args...))
}

func LogClientError(ctx context.Context, client *models.APIClient, err error) {
	frame := GetCallerFunctionName()
	writeLog(ctx, client, "E: "+frame+": "+err.Error())
}
