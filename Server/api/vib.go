package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/protobuf/types/known/emptypb"
)

func sendControlMessage(action string) {
	json, err := json.Marshal(map[string]string{"action": action})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://vib.local.mooody.me/action", "application/json", bytes.NewReader(json))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func (apiServer *MoodyAPIServer) StartViberator(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	sendControlMessage("start")
	return &emptypb.Empty{}, nil
}

func (apiServer *MoodyAPIServer) StopViberator(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	sendControlMessage("stop")
	return &emptypb.Empty{}, nil
}
