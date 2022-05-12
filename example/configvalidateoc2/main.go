/*
Configures a BGP neighbor using a BGP OpenConfig model template.
It immediately starts listening for BGP neighbor state from OpenConfig model.
*/

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	xr "github.com/nleiva/xrgrpc"
	"github.com/nleiva/xrgrpc/proto/telemetry"
	"google.golang.org/protobuf/proto"
)

// Colors, just for fun.
const (
	// blue   = "\x1b[34;1m"
	white  = "\x1b[0m"
	red    = "\x1b[31;1m"
	green  = "\x1b[32;1m"
	yellow = "\x1b[33;1m"
)

// NeighborConfig uses asplain notation for AS numbers (RFC5396)
type NeighborConfig struct {
	LocalAs         uint32
	PeerAs          uint32
	Description     string
	NeighborAddress string
	LocalAddress    string
}

func main() {
	// YANG template; defaults to "../input/template/oc-bgp.json"
	templ := flag.String("bt", "../input/template/oc-bgp.json", "BGP Config Template")
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

	neighbor := &NeighborConfig{
		LocalAs:         64512,
		PeerAs:          64512,
		Description:     "iBGP session",
		NeighborAddress: "2001:db8:cafe::2",
		//NeighborAddress: "2001:f00:bb::2",
	}

	// Read the template file
	t, err := template.ParseFiles(*templ)
	if err != nil {
		log.Fatalf("could not read the template file:  %v", err)
	}

	// 'buf' is an io.Writter to capture the template execution output
	buf := new(bytes.Buffer)
	err = t.Execute(buf, neighbor)
	if err != nil {
		log.Fatalf("could not execute the template: %v", err)
	}

	conn, ctx, err := xr.Connect(*router)
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", router.Host, err)
	}
	defer conn.Close()

	// Apply the template+parameters to the target
	_, err = xr.MergeConfig(ctx, conn, buf.String(), id)
	if err != nil {
		log.Fatalf("failed to config %s: %v\n", router.Host, err)
	} else {
		fmt.Printf("\n1)\nConfig merged on %s -> Request ID: %v\n", router.Host, id)
	}

	id++
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c := make(chan os.Signal, 1)
	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// subscription
	p := "BGP-OC"
	// encoding gpbkv
	var e int64 = 3
	ch, ech, err := xr.GetSubscription(ctx, conn, p, id, e)
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
	fmt.Printf("\n2)\nTelemetry from %s ->\n\n", router.Host)

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
	var color string
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
