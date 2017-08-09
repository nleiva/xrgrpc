// Big TODO
package xrgrpc

import (
	"net"
	"os"
	"strings"
	"testing"
)

// We are mainly validating the target config file for now
func TestDecodeJSONConfig(t *testing.T) {
	targets := NewDevices()
	location := "example/input/config.json"
	err := DecodeJSONConfig(targets, location)
	if err != nil {
		t.Fatalf("could not open the file: %s; %v", location, err)
	}
	for _, router := range targets.Routers {
		if router.User == "" {
			t.Fatalf("invalid username for %v", router.Host)
		}
		if router.Password == "" {
			t.Fatalf("invalid password for %v", router.Host)
		}
		_, err := net.ResolveTCPAddr("tcp", router.Host)
		if err != nil {
			t.Fatalf("%v is not a valid target: %v", router.Host, err)
		}
		creds := strings.Replace(router.Creds, "..", "example", 1)
		if _, err := os.Stat(creds); os.IsNotExist(err) {
			t.Fatalf("%v cert file not found for %v", router.Host, err)
		}
		if !(router.Timeout > 0) {
			t.Fatalf("invalid timeout for %v", router.Host)
		}
	}

}
