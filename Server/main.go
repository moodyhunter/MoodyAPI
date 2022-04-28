package main

import (
	"flag"
	"log"
	"net"

	"gopkg.in/ini.v1"

	"api.mooody.me/api"
	"api.mooody.me/db"
	"api.mooody.me/messaging"
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

	// Setup database
	dbConfigSection := config.Section("Database")
	DBAddress := dbConfigSection.Key("Address").MustString("localhost")
	DBDatabase := dbConfigSection.Key("Database").MustString("moodyapi")
	DBUser := dbConfigSection.Key("Username").MustString("moodyapi")
	DBPassword := dbConfigSection.Key("Password").String()
	db.SetupDBConnection(DBAddress, DBDatabase, DBUser, DBPassword)

	// Setup gRPC Server
	grpcServerAddress := config.Section("gRPC").Key("ListenAddress").MustString("127.0.0.1:1920")
	listener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		log.Fatalf("Failed to start API Server, %s", err)
	}
	_, grpcServer := api.CreateServer()
	log.Printf("gRPC server started on: %s", grpcServerAddress)

	// Setup Telegram Bot
	TgBotEnabled := config.Section("Telegram").Key("Enabled").MustBool(false)
	TgBotToken := config.Section("Telegram").Key("BotToken").MustString("")
	TgTargetChatId := config.Section("Telegram").Key("TargetGroup").MustInt64(0)

	messaging, err := messaging.NewTelegramMessaging(TgBotEnabled, TgBotToken, TgTargetChatId)
	if err != nil {
		log.Fatal(err)
	}
	messaging.SendMessage("Moody API is up and running.")

	grpcServer.Serve(listener)
}
