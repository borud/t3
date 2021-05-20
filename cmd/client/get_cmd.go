package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/borud/t3/pkg/apipb"
)

type getCmd struct{}

func (r *getCmd) Execute(args []string) error {
	client := newClient()

	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout*time.Second)
	defer cancel()

	var maps []*apipb.Map

	for _, idString := range args {
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Fatalf("error parsing id %s: %v", idString, err)
		}

		m, err := client.GetMap(ctx, &apipb.GetMapRequest{Id: id})
		if err != nil {
			log.Printf("error getting id %d: %v", id, err)
			continue
		}

		maps = append(maps, m)
	}

	json, err := json.MarshalIndent(maps, "", "  ")
	if err != nil {
		log.Fatalf("error marshalling results into JSON: %v", err)
	}
	fmt.Printf("%s\n", json)

	return nil
}
