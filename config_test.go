package xrgrpc_test

import (
	"net"
	"os"
	"strings"
	"testing"

	xr "github.com/nleiva/xrgrpc"
)

const (
	defaultFile  = "example/input/config.json"
	wrongFile    = "example/input/config2.json"
	wrongFileErr = "could not open the file"
)

func TestDecodeJSONConfig(t *testing.T) {
	targets := xr.NewDevices()
	tt := []struct {
		name string
		file string
		err  string
	}{
		{name: "local file", file: defaultFile},
		{name: "wrong file", file: wrongFile, err: wrongFileErr},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := xr.DecodeJSONConfig(targets, tc.file)
			if err != nil {
				if strings.Contains(err.Error(), wrongFileErr) && tc.err == wrongFileErr {
					return
				}
				t.Fatalf("could not open the file: %s; %v", tc.file, err)
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
				if !(router.Timeout > 0) {
					t.Fatalf("invalid timeout for %v", router.Host)
				}
			}
		})
	}
}
