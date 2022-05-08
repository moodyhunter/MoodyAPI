package db

import (
	"context"

	"api.mooody.me/models/dns"
)

func QueryDnsRecordWithType(hostname string, type_ string) (string, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return "", err
	}

	clientORM := dns.DNSRecordORM{}
	err = database.NewSelect().
		Model(&clientORM).
		Where("hostname = ?", hostname).
		Where("type = ?", type_).
		Scan(context.Background())

	if err != nil {
		return "", err
	}

	return clientORM.Ip, nil
}
