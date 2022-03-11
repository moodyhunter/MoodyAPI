package main

import (
	"flag"
	"log"
	"time"

	"gopkg.in/ini.v1"

	"api.mooody.me/api"
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

	config, err := ini.Load(ConfigPath)

	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	common.APISecret = config.Section("Global").Key("APISecret").String()

	if len(common.APISecret) == 0 {
		log.Fatalln("Must set APISecret in Global section.")
	}

	grpc_addr := config.Section("gRPC").Key("ListenAddress").MustString("127.0.0.1:1920")

	api.StartAPIServer(grpc_addr)
	for {
		time.Sleep(10 * time.Second)
	}
}
