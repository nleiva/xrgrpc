/*
gRPC Client
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	xr "github.com/nleiva/xrgrpc"
)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("This process took %s\n", elapsed)
}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	// CLI config to apply; defaults to "interface lo1 desc test"
	cli := flag.String("cli", "interface lo1 desc test", "Config to apply")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")
	flag.Parse()
	// Determine the ID for the transaction.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)

	// Define target devices from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("Could not read the config: %v", err)
	}

	cs := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(targets.Routers))

	for _, router := range targets.Routers {
		go func(d xr.CiscoGrpcClient, c string, i int64) {
			defer wg.Done()
			// Setup a connection to the target
			conn, ctx, err := xr.Connect(d)
			if err != nil {
				cs <- fmt.Sprintf("Could not setup a client connection to %s, %v\n", d.Host, err)
				return
			}
			defer conn.Close()

			// Apply 'cli' config to target
			err = xr.CLIConfig(ctx, conn, c, i)
			if err != nil {
				cs <- fmt.Sprintf("Failed to config %s, %v\n", d.Host, err)
				return
			}
			cs <- fmt.Sprintf("\nConfig applied to %s\n\n", d.Host)
			return

		}(router, *cli, id)
	}

	go func() {
		for v := range cs {
			fmt.Println(v)
		}
	}()
	wg.Wait()

}
