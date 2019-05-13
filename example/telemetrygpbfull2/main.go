/*
gRPC Client
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	xr "github.com/nleiva/xrgrpc"
	"github.com/nleiva/xrgrpc/proto/telemetry"
	lldp "github.com/nleiva/xrgrpc/proto/telemetry/lldp65x"
	"github.com/pkg/errors"
)

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func main() {
	// Subs options; LLDP, we will add some more
	p := flag.String("subs", "LLDP", "Telemetry Subscription")
	// Encoding option; defaults to GPB (only one supported in this example)
	enc := flag.String("enc", "gpb", "Encoding: 'json', 'gpb' or 'gpbkv'")
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
	id := r.Int63n(1000)

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("could not read the config: %v\n", err)
	}

	// Setup a connection to the target. 'd' is the index of the router
	// in the config file.
	d := 1
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

	for tele := range ch {
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(tele, message)
		if err != nil {
			log.Fatalf("could not unmarshall the message: %v\n", err)
		}
		e := message.GetEncodingPath()
		t := message.GetMsgTimestamp()
		fmt.Printf("Time %v, Path: %v\n", t, e)

		encPath := EncodingPath(e)
		for _, row := range message.GetDataGpb().GetRow() {
			// Keys
			msgStr := reflect.New(path2msg[encPath][protoKeys])
			output, err := decode(row.GetKeys(), msgStr.Interface().(proto.Message))
			if err != nil {
				log.Fatalf("could decode Keys: %v\n", err)
			}
			fmt.Printf("Decoded Keys:\n%v\n", output)
			// Content
			msgStr = reflect.New(path2msg[encPath][protoContent])
			output, err = decode(row.GetContent(), msgStr.Interface().(proto.Message))
			if err != nil {
				log.Fatalf("could not decode Content: %v\n", err)
			}
			fmt.Printf("Decoded Content:\n%v\n", output)
		}
	}
}

func decode(bk []byte, m proto.Message) (string, error) {
	err := proto.Unmarshal(bk, m)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshall the message")
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", errors.Wrap(err, "could not marshall into JSON")
	}
	b, err = prettyprint(b)
	if err != nil {
		return "", errors.Wrap(err, "could not pretty-print the message")
	}
	return string(b), err
}

const (
	protoKeys    = 0
	protoContent = 1
)

// EncodingPath is the path as reported by GetEncodingPath()
type EncodingPath string

var path2msg = map[EncodingPath][]reflect.Type{
	"Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail": []reflect.Type{
		reflect.TypeOf((*lldp.LldpNeighborEntry_KEYS)(nil)).Elem(),
		reflect.TypeOf((*lldp.LldpNeighborEntry)(nil)).Elem()},
}
