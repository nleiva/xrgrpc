/*
gRPC Client
*/

package main

import (
	"flag"
	"fmt"
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

	// CLI config to apply; defaults to "interface lo1 desc test"
	cli := flag.String("cli", "interface lo1 desc test", "Config to apply")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	// YANG path arguments; defaults to "yangpaths.json"

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
	conn, err := xr.Connect(targets.Routers[0])
	if err != nil {
		log.Fatalf("Could not setup a client connection to %s, %v", targets.Routers[0].Host, err)
	}
	defer conn.Close()

	// Apply 'cli' config to target
	err = xr.CLIConfig(conn, *cli, id)
	if err != nil {
		log.Fatalf("Failed to config %s, %v", targets.Routers[0].Host, err)
	} else {
		fmt.Printf("\nConfig applied to %s\n\n", targets.Routers[0].Host)
	}
}
