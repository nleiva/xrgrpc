package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	// YANG path arguments; defaults to "yangpaths.json"
	file := flag.String("file", "../input/base.cfg", "Config to apply on target")
	flag.Parse()
	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("could not read the config: %v\n", err)
	}

	// Setup a connection to the target. 'd' is the index of the router
	// in the config file
	d := 1
	conn, ctx, err := xr.Connect(targets.Routers[d])
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", targets.Routers[d].Host, err)
	}
	defer conn.Close()

	// Get config to apply on target
	base, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("could not read file: %v: %v\n", *file, err)
	}

	err = xr.CommitReplace(ctx, conn, string(base), "", id)
	if err != nil {
		log.Fatalf("could not apply the config to %s, %v", targets.Routers[d].Host, err)
	}
	fmt.Printf("\nconfig replaced on %s\n", targets.Routers[d].Host)
}
