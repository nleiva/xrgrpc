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

func main() {
	// CLI to issue; defaults to "show grpc status"
	cli := flag.String("cli", "show grpc status", "Command to execute")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	// YANG path arguments; defaults to "yangpaths.json"

	flag.Parse()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(1000)

	// Define target parameters from the configuration file
	target := xr.NewCiscoGrpcClient()
	err := xr.DecodeJSONConfig(target, *cfg)
	if err != nil {
		log.Fatalf("Could not read the config: %v", err)
	}

	// Setup a connection to the target
	conn, err := xr.Connect(*target)
	if err != nil {
		log.Fatalf("Could not setup a client connection to the target: %v", err)
	}
	defer conn.Close()

	// Apply 'cli' config to target
	err = xr.CLIConfig(conn, *cli, id)
	if err != nil {
		fmt.Printf("Failed to config the device: %v\n", err)
	} else {
		fmt.Println("Config Applied")
	}
}
