//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate protoc --go_out=camapi --go_opt=paths=source_relative --go-grpc_out=camapi --go-grpc_opt=paths=source_relative -I=.. CameraAPI.proto

package main
