/*
gRPC Client
*/

package main

import (
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
		Host:     "[2001:420:2cff:1204::5502:1]:56500",
		Creds:    "keys/ems.pem",
		Options:  "ems.cisco.com",
		Timeout:  5})

	if err != nil {
		log.Fatalf("Could not setup a client connection to the target: %v", err)
	}
	defer conn.Close()

	cli := "show run"
	id := rand.Int63n(1000)
	output, err := xr.ShowCmdTextOutput(conn, cli, id)
	fmt.Print(output)

}
