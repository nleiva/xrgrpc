/*
gRPC Client
*/

package main

import (
	"flag"
	"log"
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

	// IPv6 prefix to setup; defaults to "2001:db8::/32"
	pfx := flag.String("pfx", "2001:db8::/32", "IPv6 prefix to setup")
	// IPv6 next-hop to setup; defaults to "2001:db8:cafe::1"
	nh := flag.String("nh", "2001:db8:cafe::1", "IPv6 next-hop to setup")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	flag.Parse()

	// Admin Distance
	var admdis uint32 = 2

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("Could not read the config: %v\n", err)
	}

	// Setup a connection to the target. 'd' is the index of the router
	// in the config file
	d := 1
	conn, _, err := xr.Connect(targets.Routers[d])
	if err != nil {
		log.Fatalf("Could not setup a client connection to %s, %v", targets.Routers[d].Host, err)
	}
	defer conn.Close()

	// CSCva95005: Return SL_NOT_CONNECTED when the init session is killed from the Client.
	err = xr.ClientInit(conn)
	if err != nil {
		log.Fatalf("Failed to initialize connection to %s, %v", targets.Routers[d].Host, err)
	}

	// VRF Register Operation (= 1),
	err = xr.VRFOperation(conn, 1, admdis)
	if err != nil {
		log.Fatalf("Failed to register the VRF Operation on %s, %v", targets.Routers[d].Host, err)
	}
	// VRF EOF Operation (= 3),
	err = xr.VRFOperation(conn, 3, admdis)
	if err != nil {
		log.Fatalf("Failed to send VRF Operation EOF to %s, %v", targets.Routers[d].Host, err)
	}
	// Route Add Operation (= 1),
	err = xr.SetRoute(conn, 1, *pfx, admdis, *nh)
	if err != nil {
		log.Fatalf("Failed to set Route on %s, %v", targets.Routers[d].Host, err)
	}

}
