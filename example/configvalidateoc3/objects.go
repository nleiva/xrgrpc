package main

import (
	"bytes"
	"fmt"
	"html/template"
)

// Colors, just for fun.
const (
	blue   = "\x1b[34;1m"
	white  = "\x1b[0m"
	red    = "\x1b[31;1m"
	green  = "\x1b[32;1m"
	yellow = "\x1b[33;1m"
)

// NeighborConfig describes the bare minimum configuration to define a
// BGP Neighbor using asplain notation for AS numbers (RFC5396)
type NeighborConfig struct {
	LocalAs         uint32
	PeerAs          uint32
	Description     string
	NeighborAddress string
	LocalAddress    string
}

/*
 Borrowing some definitions from:
 https://github.com/openconfig/public/blob/master/release/models/telemetry/openconfig-telemetry.yang
*/

// SensorGroup is a list of telemetry sensory groups on the local system,
// where a sensor grouping represents a resuable grouping of multiple
// paths and exclude filters. Single path in this example!
type SensorGroup struct {
	SensorGroupID string
	Path          string
}

// Subscription holds information relating to persistent telemetry subscriptions.
// A persistent telemetry subscription is configued locally on the device through
// configuration, and is persistent across device restarts or other redundancy changes
type Subscription struct {
	SubscriptionID string
	SampleInterval uint64
	SensorGroup
}

// TelemetryConfig is Top level configuration and state for the device's
// telemetry system
type TelemetryConfig struct {
	Subscription
}

// templateProcess reads the template from a file and apply the parameters
// provided through an interface.
func templateProcess(file string, p interface{}) (out string, err error) {
	// Read the template file
	t, err := template.ParseFiles(file)
	if err != nil {
		return out, fmt.Errorf("could not read the template file:  %v", err)
	}
	// 'buf' is an io.Writter to capture the template execution output
	buf := new(bytes.Buffer)
	err = t.Execute(buf, p)
	if err != nil {
		return out, fmt.Errorf("could not execute the template: %v", err)
	}
	return buf.String(), nil
}
