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
	ypath := flag.String("ypath", "../input/yangconfig.json", "YANG path arguments")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")

	flag.Parse()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("Could not read the config: %v\n", err)
	}

	// Setup a connection to the target
	conn, ctx, err := xr.Connect(targets.Routers[0])
	if err != nil {
		log.Fatalf("Could not setup a client connection to %s, %v", targets.Routers[0].Host, err)
	}
	defer conn.Close()

	// Get YANG config file
	js, err := ioutil.ReadFile(*ypath)
	if err != nil {
		log.Fatalf("Could not read file: %v: %v\n", *ypath, err)
	}

	// Apply 'js' config to target
	ri, err := xr.MergeConfig(ctx, conn, string(js), id)
	if err != nil {
		log.Fatalf("Failed to config %s: %v\n", targets.Routers[0].Host, err)
	} else {
		fmt.Printf("\nConfig merged on %s -> Request ID: %v, Response ID: %v\n\n", targets.Routers[0].Host, id, ri)
	}
}
