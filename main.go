package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"

	"mooody.me/moodyapi/command/camera"
	"mooody.me/moodyapi/command/ddns"
	"mooody.me/moodyapi/command/ping"
	"mooody.me/moodyapi/common"
)

func main() {
	args := flag.Args()
	var ConfigPath string
	if len(args) < 1 {
		ConfigPath = "/etc/moodyapi/moodyapi.ini"
	} else {
		ConfigPath = args[0]
	}

	var err error
	common.ConfigFile, err = ini.Load(ConfigPath)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	listen_addr := common.ConfigFile.Section("Global").Key("ListenAddress").MustString("127.0.0.1:1919")
	secret_path := common.ConfigFile.Section("Global").Key("SecretPath").String()

	log.Println("MoodyAPI listen at:", listen_addr)
	if len(secret_path) == 0 {
		log.Fatalln("Must set SecretPath in Global section.")
	}

	r := gin.Default()

	// Server Ping
	r.GET("/"+secret_path+"/ping", ping.HandlePing)

	// Dynamic DNS Processing
	r.GET("/"+secret_path+"/ddns/", ddns.List)
	r.GET("/"+secret_path+"/ddns/:ddns", ddns.Get)
	r.DELETE("/"+secret_path+"/ddns/:ddns", ddns.Delete)
	r.POST("/"+secret_path+"/ddns/:ddns/update", ddns.Update)

	// Camera motion notification
	r.POST("/"+secret_path+"/camera/:id/notify", camera.TriggerPushNotification)

	r.Run(listen_addr)

}
