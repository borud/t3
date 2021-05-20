package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/borud/t3/pkg/apipb"
)

type addCmd struct {
	Data     string `long:"data" description:"data to enter into the map server"`
	Filename string `long:"file" description:"filename of map to enter into server"`
}

func (r *addCmd) Execute([]string) error {
	m := &apipb.Map{
		Timestamp: uint64(time.Now().UnixNano() / time.Hour.Milliseconds()),
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

	addResponse, err := client.AddMap(ctx, m)
	if err != nil {
		log.Fatalf("error adding data: %v", err)
	}

	fmt.Printf("added %d\n", addResponse.Id)

	return nil
}
