package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"

	"api.mooody.me/command/camera"
	"api.mooody.me/command/ddns"
	"api.mooody.me/command/ping"
	"api.mooody.me/common"
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

	prefix := "/" + secret_path

	// Server Ping
	r.GET(prefix+"/ping", ping.HandlePing)

	// Dynamic DNS Processing
	r.GET(prefix+"/ddns/", ddns.List)
	r.GET(prefix+"/ddns/:ddns", ddns.Get)
	r.DELETE(prefix+"/ddns/:ddns", ddns.Delete)
	r.POST(prefix+"/ddns/:ddns/update", ddns.Update)

	// Camera motion notification
	r.POST(prefix+"/camera/notify", camera.TriggerPushNotification)

	// Camera operations
	r.GET(prefix+"/camera/state", camera.StartCamera)
	r.POST(prefix+"/camera/start", camera.StartCamera)
	r.POST(prefix+"/camera/stop", camera.StopCamera)
	r.GET(prefix+"/camera/videolist", camera.StopCamera)

	r.Run(listen_addr)
}
