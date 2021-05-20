package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/borud/t3/pkg/apipb"
	"github.com/borud/t3/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var svc *service.Service
var waitForGRPCServerToStart sync.WaitGroup

const (
	// The address to which our gRPC server listens
	grpcListenAddress = ":4455"

	// The addres to which our http server listens
	httpListenAddress = ":5566"
)

func main() {
	// Create the service
	svc = service.New()

	waitForGRPCServerToStart.Add(1)

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
	//grpcServer := grpc.NewServer(grpc.UnaryInterceptor(service.LoggingServerInterceptor))
	grpcServer := grpc.NewServer()

	// Register the service with the server
	apipb.RegisterMapsServer(grpcServer, svc)

	// Start the server.  This call blocks until the server terminates.
	log.Printf("Starting gRPC server on %s", grpcListenAddress)

	// Signal that it is safe to connect to the gRPC server now
	waitForGRPCServerToStart.Done()

	// Start the gRPC server.  This will block.
	grpcServer.Serve(listenSocket)
}

// Start the REST service
func restService() {
	// Before trying to connect to the server we need to wait until it has started.
	waitForGRPCServerToStart.Wait()

	// create a client that talks to the local gRPC server
	conn, err := grpc.Dial(grpcListenAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error dialing gRPC server: %v", err)
	}
	client := apipb.NewMapsClient(conn)

	// Create new grpc-gateway mux
	mux := runtime.NewServeMux()

	// We need a context, but you can ignore this for now
	ctx := context.Background()

	// Register the handler for Maps
	err = apipb.RegisterMapsHandlerClient(ctx, mux, client)
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
