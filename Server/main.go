package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/ini.v1"

	"api.mooody.me/api"
	"api.mooody.me/db"
	"api.mooody.me/dns_server"
	"api.mooody.me/messaging"
)

var (
	tgBot     *messaging.TelegramBot
	dnsServer *dns_server.DnsServer
)

func main() {
	args := flag.Args()
	log.Printf("Args: %v", args)
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

	log.Printf("Loading config from %s", ConfigPath)

	// Setup database
	{
		dbSection := config.Section("Database")
		DBAddress := dbSection.Key("Address").MustString("localhost")
		DBDatabase := dbSection.Key("Database").MustString("moodyapi")
		DBUser := dbSection.Key("Username").MustString("moodyapi")
		DBPassword := dbSection.Key("Password").String()
		db.SetupConnection(DBAddress, DBDatabase, DBUser, DBPassword)
	}

	// Setup gRPC Server

	grpcSection := config.Section("gRPC")
	apiAddress := grpcSection.Key("ListenAddress").MustString("127.0.0.1:1920")
	api.CreateServer(apiAddress)
	go api.APIServer.Serve()

	// Setup Telegram Bot
	tgSection := config.Section("Telegram")
	TgBotIsEnabled := tgSection.Key("Enabled").MustBool(false)
	TgBotApiToken := tgSection.Key("BotToken").MustString("")
	TgBotSafeChatId := tgSection.Key("TargetGroup").MustInt64(0)
	TgBotSafeUserId := tgSection.Key("TargetUser").MustInt64(0)

	if TgBotIsEnabled {
		tgBot = messaging.NewTelegramBot(TgBotApiToken, TgBotSafeChatId, TgBotSafeUserId)
		log.Printf("Telegram bot is enabled")
		go api.APIServer.SubscribeNotificationInternal(tgBot.SendNotification)
		log.Printf("Telegram bot is subscribed to notifications")
		go tgBot.SendMessage("我起来了")
		go tgBot.ServeBotCommand()
	}

	// Setup DNS Server
	dnsSection := config.Section("DNS")
	dnsServerIsEnabled := dnsSection.Key("Enabled").MustBool(false)
	dnsServerAddress := dnsSection.Key("ListenAddress").MustString("127.0.0.1:53")
	dnsServerBaseDomain := dnsSection.Key("BaseDomain").MustString("local.mooody.me.")
	dnsServerTtl := (uint32)(dnsSection.Key("TTL").MustInt(60))

	if dnsServerIsEnabled {
		dnsServer = dns_server.NewDnsServer(dnsServerAddress, "udp", dnsServerBaseDomain, dnsServerTtl)
		go func() {

			log.Println("starting DNS server")
			err := dnsServer.Server.ListenAndServe()
			if err != nil {
				if TgBotIsEnabled {
					tgBot.SendMessage("我起不来 (DNS)")
				}
				log.Fatal(err)
			}
		}()
	}

	log.Printf("MoodyAPI is now ready")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh

	log.Printf("signal %d received, shutting down...", sig)

	api.APIServer.Stop()

	if dnsServerIsEnabled {
		dnsServer.Close()
	}

	if TgBotIsEnabled {
		tgBot.SendMessage("我走了")
		tgBot.Close()
	}

	db.ShutdownConnection()
	log.Println("bye")
}
