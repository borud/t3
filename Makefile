all: gen test vet build

build: client server

client:
	@cd cmd/client && go build -o ../../bin/client

server:
	@cd cmd/server && go build -o ../../bin/server

gen:
	@buf generate --path proto/*

vet:
	@go vet ./...

test:
	@go test ./...

init:
	@go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/bufbuild/buf/cmd/buf
