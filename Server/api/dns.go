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

	err = db.CreateDNSRecord(request.Record.Name, request.Record.Type, request.Record.Ip)
	if err != nil {
		return &dns.CreateDNSRecordResponse{Success: false}, err
	}

	return &dns.CreateDNSRecordResponse{Success: true}, nil
}

func (s *MoodyAPIServer) DeleteDNSRecord(ctx context.Context, request *dns.DeleteDNSRecordRequest) (*dns.DeleteDNSRecordResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `deletes dns record`)

	err = db.DeleteDNSRecord(request.Name, request.Type)
	if err != nil {
		return &dns.DeleteDNSRecordResponse{Success: false}, err
	}

	return &dns.DeleteDNSRecordResponse{Success: true}, nil
}

func (s *MoodyAPIServer) ListDNSRecords(ctx context.Context, request *dns.ListDNSRecordsRequest) (*dns.ListDNSRecordsResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `lists dns record`)

	data, err := db.ListDNSRecords()
	if err != nil {
		return &dns.ListDNSRecordsResponse{Entries: []*dns.DNSRecord{}}, err
	}

	return &dns.ListDNSRecordsResponse{Entries: data}, nil
}

func (s *MoodyAPIServer) UpdateDNSRecord(ctx context.Context, request *dns.UpdateDNSRecordRequest) (*dns.UpdateDNSRecordResponse, error) {
	client, err := db.AuthenticateClient(ctx, request.Auth, false)
	if err != nil {
		return nil, err
	}

	common.LogClientOperation(ctx, client, `updates dns record`)

	err = db.UpdateDNSRecord(request.Record.Name, request.Record.Type, request.Record.Ip)
	if err != nil {
		return &dns.UpdateDNSRecordResponse{Success: false}, err
	}

	return &dns.UpdateDNSRecordResponse{Success: true}, nil
}
