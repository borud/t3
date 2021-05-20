package main

import (
	"log"
	"time"

	"github.com/borud/t3/pkg/apipb"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
)

var opt struct {
	GRPCAddress string        `long:"grpc-addr" default:"127.0.0.1:4455" description:"gRPC address"`
	Verbose     bool          `short:"v" long:"verbose" description:"verbose output"`
	Timeout     time.Duration `long:"timeout" default:"5s" description:"timeout"`

	Add    addCmd    `command:"add" description:"add map entry"`
	Get    getCmd    `command:"get" description:"remove map entry"`
	List   listCmd   `command:"list" description:"list map entries"`
	Update updateCmd `command:"update" description:"list map entries"`
	Remove removeCmd `command:"remove" description:"remove map entry"`
}

func main() {
	p := flags.NewParser(&opt, flags.Default)
	p.Parse()
}

// newClient creates a new client given the parmeters it finds in the options.
func newClient() apipb.MapsClient {
	conn, err := grpc.Dial(opt.GRPCAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect to %s: %v", opt.GRPCAddress, err)
	}

	return apipb.NewMapsClient(conn)
}
