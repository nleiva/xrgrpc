/*
gRPC Client
*/

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	proto "github.com/golang/protobuf/proto"
	xr "github.com/nleiva/xrgrpc"
	"github.com/nleiva/xrgrpc/proto/telemetry"
)

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

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
	e := mape[*enc]
	if e == 0 {
		log.Fatalf("Encoding option '%v' not supported", *enc)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("Could not read the config: %v\n", err)
	}

	// Setup a connection to the target
	conn, err := xr.Connect(targets.Routers[0])
	if err != nil {
		log.Fatalf("Could not setup a client connection to %s, %v", targets.Routers[0].Host, err)
	}
	defer conn.Close()

	ch, err := xr.GetSubscription(conn, *p, id, e)
	if err != nil {
		log.Fatalf("Could not setup Telemetry Subscription: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Exit the example")
		// panic("Show me the stack")
		os.Exit(0)
	}()

	for {
		tele := <-ch
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(tele, message)
		if err != nil {
			log.Fatalf("Could not unmarshall the message: %v\n", err)
		}
		fmt.Printf("Time %v, Path: %v\n", message.GetMsgTimestamp(), message.GetEncodingPath())

		b, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Could not marshall into JSON: %v\n", err)
		}
		bjs, err := prettyprint(b)
		if err != nil {
			log.Fatalf("Could not pretty-print the message: %v\n", err)
		}
		fmt.Println(string(bjs))
	}
}
