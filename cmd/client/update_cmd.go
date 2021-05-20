package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/borud/t3/pkg/apipb"
)

type updateCmd struct {
	ID       uint64 `long:"id" description:"ID of map to update" required:"true"`
	Data     string `long:"data" description:"data to enter into the map server"`
	Filename string `long:"file" description:"filename of map to enter into server"`
}

func (r *updateCmd) Execute(args []string) error {
	m := &apipb.Map{
		Id: r.ID,
	}

	if r.Data != "" {
		m.Data = []byte(r.Data)
	}

	if r.Filename != "" {
		var err error
		m.Data, err = ioutil.ReadFile(r.Filename)
		if err != nil {
			log.Fatalf("error reading file %s: %v", r.Filename, err)
		}
	}

	if len(m.Data) == 0 {
		log.Fatal("cannot proceed with empty data")
	}

	client := newClient()

	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout*time.Second)
	defer cancel()

	_, err := client.Update(ctx, m)
	if err != nil {
		log.Fatalf("error updating data: %v", err)
	}

	fmt.Printf("updated %d\n", m.Id)

	return nil
}
