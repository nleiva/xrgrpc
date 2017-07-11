/*
gRPC Client
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"

	xr "github.com/nleiva/xrgrpc/client"
)

func main() {
	// Encoding option; defaults to JSON
	option := flag.String("enc", "json", "Encoding: 'json' or 'text'")
	// CLI to issue; defaults to "show grpc status"
	cli := flag.String("cli", "show grpc status", "Command to execute")
	// Config file; defaults to "router.conf"
	file := flag.String("cfg", "config.json", "Configuration file")
	flag.Parse()
	id := rand.Int63n(1000)
	output := "Empty"

	target := xr.NewCiscoGrpcClient()
	err := decodeJSONConfig(target, *file)
	if err != nil {
		log.Fatalf("Could not read the config: %v", err)
	}

	conn, err := xr.Connect(*target)
	if err != nil {
		log.Fatalf("Could not setup a client connection to the target: %v", err)
	}
	defer conn.Close()

	switch *option {
	case "text":
		output, err = xr.ShowCmdTextOutput(conn, *cli, id)
	case "json":
		output, err = xr.ShowCmdJsonOutput(conn, *cli, id)
	default:
		fmt.Printf("Don't recognize encoding: %v\n", *option)
	}
	// Add error checking
	fmt.Print(output)

}
