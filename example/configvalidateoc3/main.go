/*
1. Configures a Streaming Telemetry subscription using an OpenConfig model template.
2. Configures a BGP neighbor using an OpenConfig model template.
3. Subscribes to the Telemetry stream to learn about BGP neighbor status.
*/

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/golang/protobuf/proto"
	xr "github.com/nleiva/xrgrpc"
	"github.com/nleiva/xrgrpc/proto/telemetry"
)

func main() {
	// OpenConfig YANG templates
	bgpt := flag.String("bt", "../input/template/oc-bgp.json", "BGP Config Template")
	telet := flag.String("tt", "../input/template/oc-telemetry.json", "Telemetry Config Template")
	flag.Parse()

	// Determine the ID for first the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)

	// Manually specify target parameters.
	router, err := xr.BuildRouter(
		xr.WithUsername("cisco"),
		xr.WithPassword("cisco"),
		xr.WithHost("[2001:420:2cff:1204::5502:1]:57344"),
		xr.WithCert("../input/certificate/ems5502-1.pem"),
		xr.WithTimeout(45),
	)
	if err != nil {
		log.Fatalf("target parameters are incorrect: %s", err)
	}

	// Connect to the router
	conn, ctx, err := xr.Connect(*router)
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", router.Host, err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// Generate the Telemetry config
	tc := new(TelemetryConfig)
	tc.SubscriptionID = "BGP-OC"
	tc.SensorGroupID = "BGPNeighbor-OC"
	tc.Path = "openconfig-bgp:bgp/neighbors/neighbor/state"
	tc.SampleInterval = 1000

	teleConfig, err := templateProcess(*telet, tc)
	if err != nil {
		log.Fatalf("Telemetry, %s", err)
	}

	// Applly Telemetry config
	_, err = xr.MergeConfig(ctx, conn, teleConfig, id)
	if err != nil {
		log.Fatalf("failed to config %s: %v\n", router.Host, err)
	} else {
		fmt.Printf("\n1)\nTelemetry config applied on %s (Request ID: %v)\n", router.Host, id)
	}
	id++

	// First Pause
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Generate the BGP config
	neighbor := &NeighborConfig{
		LocalAs:         64512,
		PeerAs:          64512,
		Description:     "iBGP session",
		NeighborAddress: "2001:db8:cafe::2",
		//NeighborAddress: "2001:f00:bb::2",
	}
	bgpConfig, err := templateProcess(*bgpt, neighbor)
	if err != nil {
		log.Fatalf("BGP, %s", err)
	}

	// Apply the BGP template+parameters to the target
	_, err = xr.MergeConfig(ctx, conn, bgpConfig, id)
	if err != nil {
		log.Fatalf("failed to config %s: %v\n", router.Host, err)
	} else {
		fmt.Printf("\n2)\nBGP Config applied on %s (Request ID: %v)\n", router.Host, id)
	}

	// Second Pause
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Encoding GPBKV
	var e int64 = 3
	id++
	ch, ech, err := xr.GetSubscription(ctx, conn, tc.SubscriptionID, id, e)
	if err != nil {
		log.Fatalf("could not setup Telemetry Subscription: %v\n", err)
	}

	go func() {
		select {
		case <-c:
			fmt.Printf("\nmanually cancelled the session to %v\n\n", router.Host)
			cancel()
			return
		case <-ctx.Done():
			// Timeout: "context deadline exceeded"
			err = ctx.Err()
			fmt.Printf("\ngRPC session timed out after %v seconds: %v\n\n", router.Timeout, err.Error())
			return
		case err = <-ech:
			// Session canceled: "context canceled"
			fmt.Printf("\ngRPC session to %v failed: %v\n\n", router.Host, err.Error())
			return
		}
	}()
	fmt.Printf("\n3)\nReceiving Telemetry from %s ->\n\n", router.Host)

	for tele := range ch {
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(tele, message)
		if err != nil {
			log.Fatalf("could not unmarshall the message: %v\n", err)
		}
		ok := false
		exploreFields(message.GetDataGpbkv(), "", neighbor.NeighborAddress, &ok)
	}
}

func exploreFields(f []*telemetry.TelemetryField, indent string, peer string, ok *bool) {
	for _, field := range f {
		switch field.GetFields() {
		case nil:
			decodeKV(field, indent, peer, ok)
		default:
			exploreFields(field.GetFields(), indent+" ", peer, ok)
		}
	}
}

func decodeKV(f *telemetry.TelemetryField, indent string, peer string, ok *bool) {
	// This is a very specific scenario, just for this example.
	color := white
	switch f.GetValueByType().(type) {
	case *telemetry.TelemetryField_StringValue:
		switch f.GetName() {
		case "neighbor-address":
			addr := f.GetStringValue()
			if addr == peer {
				*ok = true
			} else {
				*ok = false
			}
		case "connection-state":
			if *ok {
				state := f.GetStringValue()
				switch state {
				case "bgp-st-active", "bgp-st-idle":
					color = red
				case "bgp-st-opensent", "bgp-st-connect", "bgp-st-openconfirm":
					color = yellow
				case "bgp-st-estab":
					color = green
				default:
					color = white
				}
				t := time.Now()
				fmt.Printf("\rNeighbor: %s, Time: %v, State: %s%s%s     ", peer, t.Format("15:04:05"), color, state, white)
				if state == "bgp-st-estab" {
					fmt.Printf("\n\nSession %sOK%s !! \n\n", green, white)
					os.Exit(0)
				}
			}
		default:
		}
	default:
	}
}
