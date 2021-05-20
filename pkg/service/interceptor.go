package service

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

// LoggingServerInterceptor logs the incoming request, calls the handler it was intended for
// logs the response it got and then returns the response and any error.
func LoggingServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// This line is logged before the request
	log.Printf("[BEFORE] server=%v method=%s req=%+v", info.Server, info.FullMethod, req)

	// Remember to call the original request
	response, err := handler(ctx, req)

	// This line is logged after we handed off the request to where it was supposed to go
	log.Printf("[ AFTER] response=%+v", response)

	// Now return the response and error
	return response, err
}
