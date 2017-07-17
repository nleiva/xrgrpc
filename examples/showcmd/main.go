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

	// Encoding option; defaults to JSON
	enc := flag.String("enc", "json", "Encoding: 'json' or 'text'")
	// CLI to issue; defaults to "show grpc status"
	cli := flag.String("cli", "show grpc status", "Command to execute")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")

	flag.Parse()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(1000)
	output := "Empty"

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

	// Return show command output based on encoding selected
	switch *enc {
	case "text":
		output, err = xr.ShowCmdTextOutput(conn, *cli, id)
	case "json":
		output, err = xr.ShowCmdJSONOutput(conn, *cli, id)
	default:
		fmt.Printf("Don't recognize encoding: %v\n", *enc)
	}
	if err != nil {
		fmt.Printf("Couldn't get the cli output: %v\n", err)
	}
	fmt.Println(output)
}
