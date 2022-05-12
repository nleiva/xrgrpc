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
	"strings"
	"time"

	xr "github.com/nleiva/xrgrpc"
	"github.com/nleiva/xrgrpc/proto/telemetry"
	"google.golang.org/protobuf/proto"
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
	// Variable for output formatting
	line := strings.Repeat("*", 90)
	sep := strings.Repeat("-", 37)

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
		LocalAs:     64512,
		PeerAs:      64512,
		Description: "iBGP session",
		//NeighborAddress: "2001:db8:cafe::2",
		NeighborAddress: "2001:f00:bb::2",
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
	// ri, err := xr.MergeConfig(ctx, conn, buf.String(), id)
	ri, err := xr.MergeConfig(ctx, conn, buf.String(), id)
	if err != nil {
		log.Fatalf("failed to config %s: %v\n", router.Host, err)
	} else {
		fmt.Println(line)
		fmt.Printf("\nconfig merged on %s -> Request ID: %v, Response ID: %v\n\n", router.Host, id, ri)
		fmt.Println(line)
	}

	// Get the BGP config from the device
	id++
	output, err := xr.GetConfig(ctx, conn, "{\"openconfig-bgp:bgp\": [null]}", id)
	if err != nil {
		log.Fatalf("could not get the config from %s, %v", router.Host, err)
	}
	fmt.Printf("BGP Config from %s\n\n", router.Host)
	fmt.Printf("\n%s\n", output)
	fmt.Println(line)

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

	fmt.Printf("\nTelemetry from %s\n\n", router.Host)

	for tele := range ch {
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(tele, message)
		if err != nil {
			log.Fatalf("could not unmarshall the message: %v\n", err)
		}
		ok := false
		ts := message.GetMsgTimestamp()
		ts64 := int64(ts * 1000000)
		fmt.Printf("%s Time %v %s\n", sep, time.Unix(0, ts64).Format("03:04:05PM"), sep)
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
	sep2 := strings.Repeat("-", 32)
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
			if *ok {
				fmt.Printf("%s Neighbor: %v %s\n", sep2, addr, sep2)
				fmt.Printf("%s%s: %s\n", indent, f.GetName(), addr)
			}
		case "connection-state", "description":
			if *ok {
				fmt.Printf("%s%s: %s\n", indent, f.GetName(), f.GetStringValue())
			}
		default:
		}
	case *telemetry.TelemetryField_BoolValue:
	case *telemetry.TelemetryField_Uint32Value:
		if *ok {
			switch f.GetName() {
			case "remote-as", "messages-received", "messages-sent", "prefixes-accepted", "prefixes-advertised":
				fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetUint32Value())
			case "connection-established-time":
				tm, err := time.ParseDuration(fmt.Sprint(f.GetUint32Value()) + "s")
				if err == nil {
					fmt.Printf("%s%s: %v\n", indent, f.GetName(), tm)
				}
			default:
			}
		}
	case *telemetry.TelemetryField_Uint64Value:
	case *telemetry.TelemetryField_BytesValue:
	case *telemetry.TelemetryField_Sint32Value:
	case *telemetry.TelemetryField_Sint64Value:
	default:
	}
}
