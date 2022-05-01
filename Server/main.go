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
	dbSection := config.Section("Database")
	DBAddress := dbSection.Key("Address").MustString("localhost")
	DBDatabase := dbSection.Key("Database").MustString("moodyapi")
	DBUser := dbSection.Key("Username").MustString("moodyapi")
	DBPassword := dbSection.Key("Password").String()
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
	tgSection := config.Section("Telegram")
	TgBotIsEnabled := tgSection.Key("Enabled").MustBool(false)
	TgBotApiToken := tgSection.Key("BotToken").MustString("")
	TgBotSafeChatId := tgSection.Key("TargetGroup").MustInt64(0)
	TgBotSafeUserId := tgSection.Key("TargetUser").MustInt64(0)

	messaging, err := messaging.NewTelegramBot(TgBotIsEnabled, TgBotApiToken, TgBotSafeChatId, TgBotSafeUserId)
	if err != nil {
		log.Fatal(err)
	}

	go messaging.SendMessage("我起来了")
	go messaging.HandleBotCommand()

	grpcServer.Serve(listener)
}
