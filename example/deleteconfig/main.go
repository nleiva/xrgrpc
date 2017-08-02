/*
gRPC Client
*/

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

	// YANG config; defaults to "yangconfig.json"
	ypath := flag.String("ypath", "../input/yangdelconfig.json", "YANG path arguments")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	flag.Parse()
	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("Could not read the config: %v\n", err)
	}

	// Setup a connection to the target. 'd' is the index of the router
	// in the config file
	d := 0
	conn, ctx, err := xr.Connect(targets.Routers[d])
	if err != nil {
		log.Fatalf("Could not setup a client connection to %s, %v", targets.Routers[d].Host, err)
	}
	defer conn.Close()

	// Get YANG config file to delete
	js, err := ioutil.ReadFile(*ypath)
	if err != nil {
		log.Fatalf("Could not read file: %v: %v\n", *ypath, err)
	}

	// Delete 'js' config on target
	ri, err := xr.DeleteConfig(ctx, conn, string(js), id)
	if err != nil {
		log.Fatalf("Failed to delete config from %s, %v", targets.Routers[d].Host, err)
	} else {
		fmt.Printf("\nConfig deleted on %s -> Request ID: %v, Response ID: %v\n\n", targets.Routers[d].Host, id, ri)
	}
}
