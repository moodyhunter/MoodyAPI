package main

import (
	"flag"
	"log"
	"time"

	"gopkg.in/ini.v1"

	"api.mooody.me/camapi"
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

	common.APISecret = common.ConfigFile.Section("Global").Key("APISecret").String()

	if len(common.APISecret) == 0 {
		log.Fatalln("Must set APISecret in Global section.")
	}

	listen_addr := common.ConfigFile.Section("DDNSAPI").Key("ListenAddress").MustString("127.0.0.1:1919")
	cameraapi_addr := common.ConfigFile.Section("CameraApi").Key("ListenAddress").MustString("127.0.0.1:1920")

	log.Println("MoodyAPI listen at:", listen_addr)

	camapi.StartAPIServer(cameraapi_addr)
	for {
		time.Sleep(10 * time.Second)
	}
}
