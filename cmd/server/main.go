package main

import (
	"log"
	"net"

	"github.com/borud/t3/pkg/apipb"
	"github.com/borud/t3/pkg/service"
	"google.golang.org/grpc"
)

// The address to which our gRPC server listens
const listenAddress = ":4455"

func main() {
	// Create TCP socket where the server listens
	listenSocket, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("error creating listen socket: %v", err)
	}

	// Create the service
	service := service.New()

	// Create the gRPC server
	grpcServer := grpc.NewServer()

	// Register the service with the server
	apipb.RegisterMapsServer(grpcServer, service)

	// Start the server.  This call blocks until the server terminates.
	log.Printf("Starting gRPC server on %s", listenAddress)
	grpcServer.Serve(listenSocket)
}
