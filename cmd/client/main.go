package main

import (
	"context"
	"log"
	"time"

	"github.com/borud/t3/pkg/apipb"
	"google.golang.org/grpc"
)

const connectAddr = "127.0.0.1:4455"

func main() {
	// First make a network connection.  We turn of transport
	// security off for now.
	conn, err := grpc.Dial(connectAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect to %s: %v", connectAddr, err)
	}

	// Then we create a client
	client := apipb.NewMapsClient(conn)

	// Create a map entry
	newMap := &apipb.Map{
		Timestamp: uint64(time.Now().UnixNano() / time.Hour.Milliseconds()),
		Data:      []byte("some SVG map"),
	}

	ctx := context.Background()

	// Add the map
	addResp, err := client.AddMap(ctx, newMap)
	if err != nil {
		log.Fatalf("error adding map: %v", err)
	}
	log.Printf("Added map with id %d", addResp.Id)

	// Get the map by id
	getResp, err := client.GetMap(ctx, &apipb.GetMapRequest{Id: addResp.Id})
	if err != nil {
		log.Fatalf("error getting map with id=%d: %v", addResp.Id, err)
	}
	log.Printf("got map: %+v", getResp)

	// Update the map
	getResp.Data = []byte("some other data that is new")

	_, err = client.Update(ctx, getResp)
	if err != nil {
		log.Fatalf("error updating map with id=%d: %v", getResp.Id, err)
	}

	// Now get the map again
	getResp, err = client.GetMap(ctx, &apipb.GetMapRequest{Id: getResp.Id})
	if err != nil {
		log.Fatalf("error getting map with id=%d: %v", addResp.Id, err)
	}
	log.Printf("got updated map: %+v", getResp)
}
