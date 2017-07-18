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

func main() {
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	// YANG path arguments; defaults to "yangpaths.json"
	ypath := flag.String("ypath", "../input/yangpaths.json", "YANG path arguments")

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

	// Get config for the YANG paths specified on 'js'
	js, err := ioutil.ReadFile(*ypath)
	if err != nil {
		log.Fatalf("Could not read file: %v: %v\n", *ypath, err)
	}
	output, err = xr.GetConfig(conn, string(js), id)
	// output, err = xr.CLIConfig(conn, "show run bgp", id)
	if err != nil {
		log.Fatalf("Could not get the config: %v\n", err)
	}
	fmt.Println(output)

}
