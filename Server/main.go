package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"

	"api.mooody.me/camapi"
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

	common.Secret = common.ConfigFile.Section("Global").Key("SecretPath").String()

	if len(common.Secret) == 0 {
		log.Fatalln("Must set SecretPath in Global section.")
	}

	listen_addr := common.ConfigFile.Section("DDNSAPI").Key("ListenAddress").MustString("127.0.0.1:1919")
	cameraapi_addr := common.ConfigFile.Section("CameraApi").Key("ListenAddress").MustString("127.0.0.1:1920")

	log.Println("MoodyAPI listen at:", listen_addr)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	prefix := "/" + common.Secret

	// Server Ping
	r.GET(prefix+"/ping", ping.HandlePing)

	// Dynamic DNS Processing
	r.GET(prefix+"/ddns/", ddns.List)
	r.GET(prefix+"/ddns/:ddns", ddns.Get)
	r.POST(prefix+"/ddns/:ddns/update", ddns.Update)

	camapi.StartAPIServer(cameraapi_addr)
	r.Run(listen_addr)
}
