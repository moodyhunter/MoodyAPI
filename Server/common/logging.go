package common

import (
	"context"
	"fmt"
	"log"

	"api.mooody.me/db"
	"api.mooody.me/models"
	"api.mooody.me/models/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getClientIdName(client *common.APIClient) (int64, string) {
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

func logWithApiClient(ctx context.Context, client *common.APIClient, prefix string, frame string, operation string) {
	cid, cname := getClientIdName(client)

	log.Printf("%s: [%s] %s", prefix, frame, operation)

	logInfo := models.OperationLog{
		ClientId:   cid,
		ClientName: cname,
		Time:       timestamppb.Now(),
		Operation:  operation,
	}

	db.LogOperation(ctx, &logInfo)
}

func LogClientOperation(ctx context.Context, client *common.APIClient, message string, args ...interface{}) {
	frame := GetCallerFunctionName()
	logWithApiClient(ctx, client, "I", frame, fmt.Sprintf(message, args...))
}

func LogClientError(ctx context.Context, client *common.APIClient, err error) {
	frame := GetCallerFunctionName()
	logWithApiClient(ctx, client, "E", frame, err.Error())
}
