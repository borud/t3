package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/borud/t3/pkg/apipb"
)

type removeCmd struct {
	ID uint64 `long:"id" description:"ID of map to remove"`
}

func (r *removeCmd) Execute(args []string) error {

	client := newClient()

	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout*time.Second)
	defer cancel()

	for _, idString := range args {
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			log.Fatalf("error parsing id %s: %v", idString, err)
		}

		_, err = client.DeleteMap(ctx, &apipb.DeleteMapRequest{Id: id})
		if err != nil {
			log.Fatalf("error removing id %d: %v", id, err)
			continue
		}

		fmt.Printf("removed %d\n", id)
	}

	return nil
}
