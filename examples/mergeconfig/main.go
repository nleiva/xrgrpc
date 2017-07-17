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
	// YANG config; defaults to "yangconfig.json"
	ypath := flag.String("ypath", "../input/yangconfig.json", "YANG path arguments")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")

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

	// Get YANG config file
	js, err := ioutil.ReadFile(*ypath)
	if err != nil {
		fmt.Printf("Could not read file: %v: %v\n", *ypath, err)
	}

	// Apply 'js' config to target
	ri, err := xr.MergeConfig(conn, string(js), id)
	if err != nil {
		fmt.Printf("Failed to config the device: %v\n", err)
	} else {
		fmt.Printf("Config Applied -> Request ID: %v, Response ID: %v\n", id, ri)
	}
}
