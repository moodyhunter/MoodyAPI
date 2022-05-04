//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
//go:generate go install github.com/moodyhunter/protoc-gen-bun@v1.3.0
//go:generate bash -c "shopt -s globstar && protoc -I=../proto ../proto/**/**.proto --go_out=models      --go_opt=paths=source_relative"
//go:generate bash -c "shopt -s globstar && protoc -I=../proto ../proto/**/**.proto --go-grpc_out=models --go-grpc_opt=paths=source_relative"
//go:generate bash -c "shopt -s globstar && protoc -I=../proto ../proto/**/**.proto --bun_out=models     --bun_opt=paths=source_relative"

package main
