package main

import (
	"flag"
	"log"
	"net"

	"gopkg.in/ini.v1"

	"api.mooody.me/api"
	"api.mooody.me/common"
	"api.mooody.me/db"
)

func main() {
	args := flag.Args()
	var ConfigPath string
	if len(args) < 1 {
		ConfigPath = "/etc/moodyapi/moodyapi.ini"
	} else {
		ConfigPath = args[0]
	}

	config, err := ini.Load(ConfigPath)

	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	common.APISecret = config.Section("").Key("APISecret").String()
	if len(common.APISecret) == 0 {
		log.Fatalln("Must set APISecret in Global section.")
	}

	dbConfigSection := config.Section("Database")
	DBAddress := dbConfigSection.Key("Address").String()
	DBDatabase := dbConfigSection.Key("Database").String()
	DBUser := dbConfigSection.Key("Username").String()
	DBPassword := dbConfigSection.Key("Password").String()
	db.SetupDBConnection(DBAddress, DBDatabase, DBUser, DBPassword)

	grpc_addr := config.Section("gRPC").Key("ListenAddress").MustString("127.0.0.1:1920")

	log.Printf("gRPC server started on: %s", grpc_addr)
	listener, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		log.Fatalf("Failed to start API Server, %s", err)
	}

	_, grpcServer := api.CreateServer()
	grpcServer.Serve(listener)
}
