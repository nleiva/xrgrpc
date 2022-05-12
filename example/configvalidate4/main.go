/*
Configures a BGP neighbor using a BGP Openconfig model template.
It inmidiately starts listening for BGP neighbor state.
*/

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	xr "github.com/nleiva/xrgrpc"
	"github.com/nleiva/xrgrpc/proto/telemetry"
	bgp "github.com/nleiva/xrgrpc/proto/telemetry/bgp4"
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

	// YANG template; defaults to "../input/template/oc-bgpv4.json"
	templ := flag.String("bt", "../input/template/oc-bgpv4.json", "BGP Config Template")
	flag.Parse()

	// Determine the ID for first the transaction.
	var id int64 = 1000

	// Manually specify target parameters.
	router1, err := xr.BuildRouter(
		xr.WithUsername("vagrant"),
		xr.WithPassword("vagrant"),
		xr.WithHost("192.0.2.2:57344"),
		xr.WithCert("ems1.pem"),
		xr.WithTimeout(60),
	)
	if err != nil {
		log.Fatalf("target parameters for router1 are incorrect: %s", err)
	}
	router2, err := xr.BuildRouter(
		xr.WithUsername("vagrant"),
		xr.WithPassword("vagrant"),
		xr.WithHost("192.0.2.3:57344"),
		xr.WithCert("ems2.pem"),
		xr.WithTimeout(15),
	)
	if err != nil {
		log.Fatalf("target parameters for router2 are incorrect: %s", err)
	}

	// Define BGP parameters for each device
	neighbor1 := &NeighborConfig{
		LocalAs:         64512,
		PeerAs:          64512,
		Description:     "iBGP session",
		NeighborAddress: "203.0.113.3",
		LocalAddress:    "203.0.113.2",
	}
	neighbor2 := &NeighborConfig{
		LocalAs:         64512,
		PeerAs:          64512,
		Description:     "iBGP session",
		NeighborAddress: "203.0.113.2",
		LocalAddress:    "203.0.113.3",
	}

	// Read the OC BGP template file
	t, err := template.ParseFiles(*templ)
	if err != nil {
		log.Fatalf("could not read the template file:  %v", err)
	}

	// 'buf' is an io.Writter to capture the template execution output for each device
	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)
	err = t.Execute(buf1, neighbor1)
	if err != nil {
		log.Fatalf("could not execute the template for router 1: %v", err)
	}
	err = t.Execute(buf2, neighbor2)
	if err != nil {
		log.Fatalf("could not execute the template for router 2: %v", err)
	}

	// Connect to the targets
	conn1, ctx1, err := xr.Connect(*router1)
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", router1.Host, err)
	}
	defer conn1.Close()
	conn2, ctx2, err := xr.Connect(*router2)
	if err != nil {
		log.Fatalf("could not setup a client connection to %s, %v", router2.Host, err)
	}
	defer conn2.Close()

	// Apply the template+parameters to the targets.
	ri, err := xr.MergeConfig(ctx1, conn1, buf1.String(), id)
	if err != nil {
		log.Fatalf("failed to config %s: %v\n", router1.Host, err)
	} else {
		fmt.Println(line)
		fmt.Printf("\nconfig merged on %s -> Request ID: %v, Response ID: %v\n\n", router1.Host, id, ri)
	}
	ri, err = xr.MergeConfig(ctx2, conn2, buf2.String(), id)
	if err != nil {
		log.Fatalf("failed to config %s: %v\n", router2.Host, err)
	} else {
		fmt.Printf("\nconfig merged on %s -> Request ID: %v, Response ID: %v\n\n", router2.Host, id, ri)
		fmt.Println(line)
	}
	// We no longer need the second connection, we focus on the first device.
	conn2.Close()

	// Get the BGP config from one of the devices
	id++
	output, err := xr.GetConfig(ctx1, conn1, "{\"openconfig-bgp:bgp\": [null]}", id)
	if err != nil {
		log.Fatalf("could not get the config from %s, %v", router1.Host, err)
	}
	fmt.Printf("\nBGP Config from %s\n\n", router1.Host)
	fmt.Printf("\n%s\n", output)
	fmt.Println(line)

	id++
	ctx1, cancel := context.WithCancel(ctx1)
	defer cancel()
	c := make(chan os.Signal, 1)
	// If no signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will. E.g.: signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// Telemetry Subscription
	p := "BGP"
	// encoding GPB
	var e int64 = 2
	ch, ech, err := xr.GetSubscription(ctx1, conn1, p, id, e)
	if err != nil {
		log.Fatalf("could not setup Telemetry Subscription: %v\n", err)
	}

	go func() {
		select {
		case <-c:
			fmt.Printf("\nmanually cancelled the session to %v\n\n", router1.Host)
			cancel()
			return
		case <-ctx1.Done():
			// Timeout: "context deadline exceeded"
			err = ctx1.Err()
			fmt.Printf("\ngRPC session timed out after %v seconds: %v\n\n", router1.Timeout, err.Error())
			return
		case err = <-ech:
			// Session canceled: "context canceled"
			fmt.Printf("\ngRPC session to %v failed: %v\n\n", router1.Host, err.Error())
			return
		}
	}()

	fmt.Printf("\ntelemetry from %s\n\n", router1.Host)

	for tele := range ch {
		message := new(telemetry.Telemetry)
		err := proto.Unmarshal(tele, message)
		if err != nil {
			log.Fatalf("could not unmarshall the message: %v\n", err)
		}
		ts := message.GetMsgTimestamp()
		ts64 := int64(ts * 1000000)
		fmt.Printf("%s Time %v %s\n", sep, time.Unix(0, ts64).Format("03:04:05PM"), sep)

		for _, row := range message.GetDataGpb().GetRow() {
			content := row.GetContent()
			nbr := new(bgp.BgpNbrBag)
			err = proto.Unmarshal(content, nbr)
			if err != nil {
				log.Fatalf("could decode Content: %v\n", err)
			}
			rasn := nbr.GetRemoteAs()
			// tn := nbr.GetConnectionEstablishedTime()
			state := nbr.GetConnectionState()
			raddr := nbr.GetConnectionRemoteAddress().GetIpv4Address()
			// Debug:
			// fmt.Printf("\n\n%v\n\n\n", hex.Dump(content))
			fmt.Printf("BGP Neighbor; IP: %v, ASN: %v, State %s \n\n", raddr, rasn, state)

		}
	}
}
