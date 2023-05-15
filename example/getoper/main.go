package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"math/rand"
	"time"

	xr "github.com/nleiva/xrgrpc"
)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("This process took %s\n", elapsed)

}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	// YANG path arguments; defaults to "yangoper.json"
	ypath := flag.String("ypath", "../input/yangoper.json", "YANG path arguments")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	flag.Parse()
	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)
	var output string

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("could not read the config: %v\n", err)
	}

	// Setup a connection to the target. 'd' is the index of the router
	// in the config file
	d := 0
	conn, ctx, err := xr.Connect(targets.Routers[d])
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", targets.Routers[d].Host, err)
	}
	defer conn.Close()

	// Get YANG config file
	js, err := os.ReadFile(*ypath)
	if err != nil {
		log.Fatalf("could not read file: %v: %v\n", *ypath, err)
	}

	output, err = xr.GetOper(ctx, conn, string(js), id)
	if err != nil {
		log.Fatalf("could not get the operation data from %s, %v", targets.Routers[d].Host, err)
	}
	fmt.Printf("\noperation data from %s\n %s\n", targets.Routers[d].Host, output)
}
