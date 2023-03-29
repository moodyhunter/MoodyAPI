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

func CreateDNSRecord(hostname string, dnstype string, ip string) error {
	err := checkDatabaseConnectivity()
	if err != nil {
		return err
	}

	clientORM := dns.DNSRecordORM{
		Name: hostname,
		Type: dnstype,
		Ip:   ip,
	}

	_, err = database.NewInsert().
		Model(&clientORM).
		Exec(context.Background())

	return err
}

func DeleteDNSRecord(hostname string, dnstype string) error {
	err := checkDatabaseConnectivity()
	if err != nil {
		return err
	}

	_, err = database.NewDelete().
		Model(&dns.DNSRecordORM{}).
		Where("hostname = ?", hostname).
		Where("type = ?", dnstype).
		Exec(context.Background())

	return err
}

func ListDNSRecords() ([]*dns.DNSRecord, error) {
	err := checkDatabaseConnectivity()
	if err != nil {
		return nil, err
	}

	clientORM := []dns.DNSRecordORM{}
	err = database.NewSelect().
		Model(&clientORM).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	var records []*dns.DNSRecord
	for _, record := range clientORM {
		records = append(records, &dns.DNSRecord{
			Name: record.Name,
			Type: record.Type,
			Ip:   record.Ip,
		})
	}

	return records, nil
}

func UpdateDNSRecord(hostname string, dnstype string, ip string) error {
	err := checkDatabaseConnectivity()
	if err != nil {
		return err
	}

	clientORM := dns.DNSRecordORM{
		Name: hostname,
		Type: dnstype,
		Ip:   ip,
	}

	_, err = database.NewUpdate().
		Model(&clientORM).
		Exec(context.Background())

	return err
}
