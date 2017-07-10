/*
gRPC Client
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"

	xr "github.com/nleiva/xrgrpc/client"
)

func main() {
	// XR self-signed certificates are issued for CN=ems.cisco.com
	conn, err := xr.Connect(xr.CiscoGrpcClient{
		User:     "cisco",
		Password: "cisco",
		Host:     "[2001:420:2cff:1204::5502:1]:57344",
		Creds:    "keys/ems.pem",
		Options:  "ems.cisco.com",
		Timeout:  5})

	if err != nil {
		log.Fatalf("Could not setup a client connection to the target: %v", err)
	}
	defer conn.Close()

	option := flag.String("e", "json", "Encoding: json or text")
	cli := flag.String("c", "show route ipv6", "Command to execute")
	flag.Parse()
	id := rand.Int63n(1000)
	output := "Empty"

	switch *option {
	case "text":
		output, err = xr.ShowCmdTextOutput(conn, *cli, id)
	case "json":
		output, err = xr.ShowCmdJsonOutput(conn, *cli, id)
	default:
		fmt.Printf("Don't recognize encoding: %v\n", option)
	}

	fmt.Print(output)

}
