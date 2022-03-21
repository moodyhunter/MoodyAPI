//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
//go:generate protoc --go_out=api/pb --go_opt=paths=source_relative --go-grpc_out=api/pb --go-grpc_opt=paths=source_relative -I=.. MoodyAPI.proto

package main
