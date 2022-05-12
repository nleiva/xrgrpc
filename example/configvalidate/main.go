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
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	xr "github.com/nleiva/xrgrpc"
	"github.com/nleiva/xrgrpc/proto/telemetry"
	bgp "github.com/nleiva/xrgrpc/proto/telemetry/bgp"
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
		xr.WithUsername("admin"),
		xr.WithPassword("C1sco12345"),
		xr.WithHost("sandbox-iosxr-1.cisco.com:57777"),
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
	input := `{"openconfig-network-instance:network-instances":
	{"network-instance":[{"name":"default","protocols":{"protocol":[{"bgp":{}}]}}]}}`

	output, err := xr.GetConfig(ctx, conn, input, id)

	if err != nil {
		log.Fatalf("could not get the config from %s, %v", router.Host, err)
	}
	fmt.Printf("\nBGP Config from %s\n\n", router.Host)
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
	p := "BGP"
	// encoding GPB
	var e int64 = 2
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
			state := nbr.GetConnectionState()
			raddr := nbr.GetConnectionRemoteAddress().Ipv6Address

			fmt.Printf("BGP Neighbor; IP: %v, ASN: %v, State %s \n\n", raddr, rasn, state)

		}
	}
}
