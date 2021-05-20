package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/borud/t3/pkg/apipb"
	"github.com/borud/t3/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var svc *service.Service

const (
	// The address to which our gRPC server listens
	grpcListenAddress = ":4455"

	// The addres to which our http server listens
	httpListenAddress = ":5566"
)

func main() {
	// Create the service
	svc = service.New()

	go grpcService()
	go restService()

	// Capture Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

// Start the gRPC service
func grpcService() {
	// Create TCP socket where the server listens
	listenSocket, err := net.Listen("tcp", grpcListenAddress)
	if err != nil {
		log.Fatalf("error creating listen socket: %v", err)
	}

	// Create the gRPC server
	grpcServer := grpc.NewServer()

	// Register the service with the server
	apipb.RegisterMapsServer(grpcServer, svc)

	// Start the server.  This call blocks until the server terminates.
	log.Printf("Starting gRPC server on %s", grpcListenAddress)
	grpcServer.Serve(listenSocket)
}

// Start the REST service
func restService() {
	// Create new grpc-gateway mux
	mux := runtime.NewServeMux()

	// We need a context, but you can ignore this for now
	ctx := context.Background()

	// Register the handler for Maps
	err := apipb.RegisterMapsHandlerServer(ctx, mux, svc)
	if err != nil {
		log.Fatalf("unable to register MapsHandlerServer: %v", err)
	}

	httpServer := &http.Server{
		Addr:    httpListenAddress,
		Handler: mux,
	}

	log.Printf("starting http server on %s", httpListenAddress)
	httpServer.ListenAndServe()
}
