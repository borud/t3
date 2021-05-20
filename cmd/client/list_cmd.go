package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

type listCmd struct{}

func (r *listCmd) Execute(args []string) error {
	client := newClient()
	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout*time.Second)
	defer cancel()

	resp, err := client.ListMaps(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error listing maps: %v", err)
	}

	if len(resp.Maps) > 0 {
		json, err := json.MarshalIndent(resp.Maps, "", "  ")
		if err != nil {
			log.Fatalf("error marshalling results into JSON: %v", err)
		}

		fmt.Printf("%s\n", json)
	}

	return nil
}
