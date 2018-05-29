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

	// Encoding option; defaults to JSON
	enc := flag.String("enc", "json", "Encoding: 'json' or 'cli'")
	// Action to issue; defaults to "ping4.json"
	act := flag.String("act", "../input/action/ping6.json", "Command to execute")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	flag.Parse()

	file, err := ioutil.ReadFile(*act)
	if err != nil {
		log.Fatalf("could not read file: %v\n", *act)
	}
	cli := string(file)

	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)
	output := "Empty"

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err = xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("could not read the config: %v\n", err)
	}

	// Setup a connection to the target. 'd' is the index of the router
	// in the config file
	d := 3
	conn, ctx, err := xr.Connect(targets.Routers[d])
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", targets.Routers[d].Host, err)
	}
	defer conn.Close()

	// Return show command output based on encoding selected
	switch *enc {
	case "json":
		output, err = xr.ActionJSON(ctx, conn, cli, id)
	//case "cli":
	//	output, err = xr.ActionCLI(ctx, conn, cli, id)
	default:
		log.Fatalf("don't recognize encoding: %v\n", *enc)
	}
	if err != nil {
		log.Fatalf("couldn't get an output: %v\n", err)
	}
	fmt.Printf("\noutput from %s\n %s\n", targets.Routers[d].Host, output)
}
