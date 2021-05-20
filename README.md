# Tutorial 3

In this tutorial we make a very simple map server.

## 01 Basic setup

Create the git project and module

    git init
    go mod init github.com/borud/t3

Create `proto` and `pkg` directory.

### Make proto file

```protobuf

syntax = "proto3";
package apipb;

import "google/protobuf/empty.proto";

option go_package = "pkg/apipb";

message Map {
 uint64 id = 1;
 uint64 timestamp = 2;
 bytes data = 3;
}

message AddMapResponse {
 uint64 id = 1;
}

message GetMapRequest {
 uint64 id = 1;
}

message DeleteMapRequest {
 uint64 id = 1;
}

service Maps {
 rpc AddMap(Map) returns (AddMapResponse);
 rpc GetMap(GetMapRequest) returns (Map);
 rpc Update(Map) returns (google.protobuf.Empty);
 rpc DeleteMap(DeleteMapRequest) returns (google.protobuf.Empty);
}
```

### Add `buf.yaml` and `buf.gen.yaml`

buf.yaml:

```yaml
version: v1beta1
name: 
build:
  roots:
    - proto
```

buf.gen.yaml:

```yaml
version: v1beta1
plugins:
  - name: go
    out: pkg/apipb
    opt: paths=source_relative
  - name: go-grpc
    out: pkg/apipb
    opt: paths=source_relative,require_unimplemented_servers=false
```

### Add `Makefile`

```makefile
gen:
	@buf generate

init:
	@go get -u github.com/grpc-ecosystem/grpc-gateway/v2\
        protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/bufbuild/buf/cmd/buf
```

## Implement service and server

### Implementing the service

First check out the generated code and look for the service interface.  Our service was named `Maps` in the proto file so the interface definition for the server should be named `MapsServer`.  

*This naming is a bit unfortunate since we tend to distinguish between service and server in our code.  The service is the business logic and the server is the thing that takes care of the plumbing.*

```go
type MapsServer interface {
 AddMap(context.Context, *Map) (*AddMapResponse, error)
 GetMap(context.Context, *GetMapRequest) (*Map, error)
 Update(context.Context, *Map) (*emptypb.Empty, error)
 DeleteMap(context.Context, *DeleteMapRequest) (*emptypb.Empty, error)
}
```

This is implemented in [`pkg/service/service.go`](pkg/service/service.go)

### Implementing the server

Then create a server which will tie the parts together.  This can be found in
[`cmd/server/main.go`](cmd/server/main.go) and shows how we tie together the parts.
