/*
gRPC Client
*/

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"time"

	xr "github.com/nleiva/xrgrpc"
)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("This process took %s\n", elapsed)

}

// NeighborConfig uses asplain notation for AS numbers (RFC5396)
type NeighborConfig struct {
	LocalAs         uint32
	PeerAs          uint32
	Description     string
	NeighborAddress string
	LocalAddress    string
}

// It uses asdot+ notation according to RFC5396
type fourByteASN struct {
	X uint16
	Y uint16
}

type bgpData struct {
	*NeighborConfig
	LocalASN fourByteASN
	PeerASN  fourByteASN
}

func splitASN(a uint32) fourByteASN {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, a)
	return fourByteASN{binary.BigEndian.Uint16(b[:2]), binary.BigEndian.Uint16(b[2:])}
}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	// YANG config template
	templ := flag.String("bt", "../input/template/bgp.json", "BGP Config Template")
	// BGP config parameters
	bgpparam := flag.String("prm", "../input/template/bgp-parameters.json", "YANG path arguments")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	flag.Parse()
	// Determine the ID for first the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("could not read the config: %v\n", err)
	}

	// Define the bgp parameters from the BGP parameters file
	p := new(NeighborConfig)
	err = xr.DecodeJSONConfig(p, *bgpparam)
	if err != nil {
		log.Fatalf("could not read the BGP parameters: %v\n", err)
	}

	// Deal with the as-xx, as-yy notation of the Cisco IOS XR YANG model.
	data := bgpData{
		p,
		splitASN(p.LocalAs),
		splitASN(p.PeerAs),
	}

	// Read the template file
	t, err := template.ParseFiles(*templ)
	if err != nil {
		log.Fatalf("could not read the template file:  %v", err)
	}

	// 'buf' is an io.Writter to capture the template execution output
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		log.Fatalf("could not execute the template: %v", err)
	}

	// Setup a connection to the target. 'd' is the index of the router
	// in the config file
	d := 0
	conn, ctx, err := xr.Connect(targets.Routers[d])
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", targets.Routers[d].Host, err)
	}
	defer conn.Close()

	// Apply the template+parameters to the target
	ri, err := xr.MergeConfig(ctx, conn, buf.String(), id)
	if err != nil {
		log.Fatalf("failed to config %s: %v\n", targets.Routers[d].Host, err)
	} else {
		fmt.Printf("\nconfig merged on %s -> Request ID: %v, Response ID: %v\n\n", targets.Routers[d].Host, id, ri)
	}

	// Get the BGP config from the device
	id++
	output, err := xr.GetConfig(ctx, conn, "{\"Cisco-IOS-XR-ipv4-bgp-cfg:bgp\": [null]}", id)
	if err != nil {
		log.Fatalf("could not get the config from %s, %v", targets.Routers[d].Host, err)
	}
	fmt.Printf("\nconfig from %s\n %s\n", targets.Routers[d].Host, output)
}
