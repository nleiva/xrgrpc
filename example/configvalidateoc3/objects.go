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

/*
Borrowing some definitions from:

- openconfig-interfaces:
 https://github.com/openconfig/public/blob/master/release/models/interfaces/openconfig-interfaces.yang

- openconfig-if-ethernet:
https://github.com/openconfig/public/blob/master/release/models/interfaces/openconfig-if-ethernet.yang

- openconfig-if-ip:
https://github.com/openconfig/public/blob/master/release/models/interfaces/openconfig-if-ip.yang

- A YANG Data Model for Interface Management (RFC 7223): https://tools.ietf.org/html/rfc7223

- Definitions of Managed Objects for the Ethernet-like Interface Types (RFC 3635):
https://tools.ietf.org/html/rfc3635

TODO: Take a look at:
https://github.com/openconfig/ygot#validating-the-struct-contents
*/

// InterfaceConfig represents the configuration of a physical interfaces
// and subinterfaces.
// At present, ethernet-like media are identified by the value ethernetCsmacd(6) of the ifType
// object in the Interfaces MIB [RFC2863]
type InterfaceConfig struct {
	Name        string
	Description string
	Enabled     bool
	Physical
	Ethernet
	SubInterface
}

// Physical defines configuration data for physical interfaces.
type Physical struct {
	Type string
	MTU  uint16
}

// Ethernet defines configuration items for Ethernet interfaces.
type Ethernet struct {
	MACAddress    string
	AutoNegotiate bool
	DuplexMode    string
	PortSpeed     string
}

// SubInterface defines data for logical interfaces associated with a given interface.
type SubInterface struct {
	Index uint32
	IPv6
}

// IPv6 defines configuration and state for IPv6 interfaces.
type IPv6 struct {
	Address      string
	PrefixLength uint8
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

// ⍻	NOT CHECK MARK (U+237B)	e28dbb
// ✅	WHITE HEAVY CHECK MARK (U+2705)	e29c85
// ✓	CHECK MARK (U+2713)	e29c93
// ✔	HEAVY CHECK MARK (U+2714) e29c94
