//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
//go:generate go install github.com/moodyhunter/protoc-gen-bun@v1.3.0
//go:generate protoc --go_out=models      --go_opt=paths=source_relative      -I=../proto MoodyAPI.proto wg/wg.proto common/common.proto privileged/privileged.proto notifications/notifications.proto
//go:generate protoc --go-grpc_out=models --go-grpc_opt=paths=source_relative -I=../proto MoodyAPI.proto wg/wg.proto common/common.proto privileged/privileged.proto notifications/notifications.proto
//go:generate protoc --bun_out=models     --bun_opt=paths=source_relative     -I=../proto MoodyAPI.proto wg/wg.proto common/common.proto privileged/privileged.proto notifications/notifications.proto

package main
