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

	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	// YANG path arguments; defaults to "yangpaths.json"
	ypath := flag.String("ypath", "../input/yangpaths.json", "YANG path arguments")
	flag.Parse()
	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)
	output := "Empty"

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

	// Get config for the YANG paths specified on 'js'
	js, err := ioutil.ReadFile(*ypath)
	if err != nil {
		log.Fatalf("Could not read file: %v: %v\n", *ypath, err)
	}
	output, err = xr.GetConfig(ctx, conn, string(js), id)
	// output, err = xr.CLIConfig(conn, "show run bgp", id)
	if err != nil {
		log.Fatalf("Could not get the config from %s, %v", targets.Routers[d].Host, err)
	}
	fmt.Printf("\nConfig from %s\n %s\n", targets.Routers[d].Host, output)
}
