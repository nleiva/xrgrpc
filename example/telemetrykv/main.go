package main

import (
	"context"
	"flag"
	"fmt"
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

func main() {
	// Subs options; LLDP, we will add some more
	p := flag.String("subs", "LLDP", "Telemetry Subscription")
	// Encoding option; defaults to GPBKV (only one supported in this example)
	enc := flag.String("enc", "gpbkv", "Encoding: 'json', 'gpb' or 'gpbkv'")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	flag.Parse()

	mape := map[string]int64{
		"gpb":   2,
		"gpbkv": 3,
		"json":  4,
	}
	e, ok := mape[*enc]
	if !ok {
		log.Fatalf("encoding option '%v' not supported", *enc)
	}

	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("could not read the config: %v\n", err)
	}

	// Setup a connection to the target. 'd' is the index of the router
	// in the config file
	d := 0
	// Adjust timeout to increase gRPC session lifespan to be able to receive
	// Streaming Telemetry data for a period of time.
	targets.Routers[d].Timeout = 20
	conn, ctx, err := xr.Connect(targets.Routers[d])
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", targets.Routers[d].Host, err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch, ech, err := xr.GetSubscription(ctx, conn, *p, id, e)
	if err != nil {
		log.Fatalf("could not setup Telemetry Subscription: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			fmt.Printf("\nmanually cancelled the session to %v\n\n", targets.Routers[d].Host)
			cancel()
			return
		case <-ctx.Done():
			// Timeout: "context deadline exceeded"
			err = ctx.Err()
			fmt.Printf("\ngRPC session timed out after %v seconds: %v\n\n", targets.Routers[d].Timeout, err.Error())
			return
		case err = <-ech:
			// Session canceled: "context canceled"
			fmt.Printf("\ngRPC session to %v failed: %v\n\n", targets.Routers[d].Host, err.Error())
			return
		}
	}()

	line := strings.Repeat("*", 90)
	for tele := range ch {
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(tele, message)
		if err != nil {
			log.Fatalf("could not unmarshall the message: %v\n", err)
		}
		ts := message.GetMsgTimestamp()
		ts64 := int64(ts * 1000000)
		fmt.Println(line)
		fmt.Printf("Time %v, Path: %v\n", time.Unix(0, ts64).Format("03:04:05PM"), message.GetEncodingPath())
		fmt.Println(line)
		exploreFields(message.GetDataGpbkv(), "")
	}

}

func exploreFields(f []*telemetry.TelemetryField, indent string) {
	for _, field := range f {
		switch field.GetFields() {
		case nil:
			decodeKV(field, indent)
		default:
			exploreFields(field.GetFields(), indent+" ")
		}
	}
}

func decodeKV(f *telemetry.TelemetryField, indent string) {
	// This is incomplete, but covers most of the cases I've seen so far.
	switch f.GetValueByType().(type) {
	case *telemetry.TelemetryField_StringValue:
		fmt.Printf("%s%s: %s\n", indent, f.GetName(), f.GetStringValue())
	case *telemetry.TelemetryField_BoolValue:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetBoolValue())
	case *telemetry.TelemetryField_Uint32Value:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetUint32Value())
	case *telemetry.TelemetryField_Uint64Value:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetUint64Value())
	case *telemetry.TelemetryField_BytesValue:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetBytesValue())
	case *telemetry.TelemetryField_Sint32Value:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetSint32Value())
	case *telemetry.TelemetryField_Sint64Value:
		fmt.Printf("%s%s: %v\n", indent, f.GetName(), f.GetSint64Value())
	default:
	}
}
