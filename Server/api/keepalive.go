package api

import (
	"context"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models"
)

func (apiServer *MoodyAPIServer) KeepAlive(request *models.KeepAliveRequest, server models.MoodyAPIService_KeepAliveServer) error {
	client, err := db.AuthenticateClient(context.Background(), request.Auth, false)
	if err != nil {
		common.LogClientError(context.Background(), client, err)
		return err
	}

	common.LogClientOperation(context.Background(), client, "connected")

	apiServer.keepAliveStream.BlockedSubscribeWithCallback(func(signal *models.KeepAliveMessage) {
		server.Send(signal)
	})

	common.LogClientOperation(context.Background(), client, "disconnected")
	return nil
}
