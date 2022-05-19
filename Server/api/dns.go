package api

import (
	"context"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models/dns"
)

func (s *MoodyAPIServer) CreateDNSRecord(ctx context.Context, request *dns.CreateDNSRecordRequest) (*dns.CreateDNSRecordResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `creates dns record`)

	// TODO
	return nil, nil
}

func (s *MoodyAPIServer) DeleteDNSRecord(ctx context.Context, request *dns.DeleteDNSRecordRequest) (*dns.DeleteDNSRecordResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `deletes dns record`)

	// TODO
	return nil, nil
}

func (s *MoodyAPIServer) ListDNSRecords(ctx context.Context, request *dns.ListDNSRecordsRequest) (*dns.ListDNSRecordsResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `lists dns record`)

	// TODO
	return nil, nil
}

func (s *MoodyAPIServer) UpdateDNSRecord(ctx context.Context, request *dns.UpdateDNSRecordRequest) (*dns.UpdateDNSRecordResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `updates dns record`)

	// TODO
	return nil, nil
}
