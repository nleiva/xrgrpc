// Package xrgrpc is a gRPC Client library for Cisco IOS XR devices. It
// exposes different RPC's to manage the device(s). The objective is
// to have a single interface to retrieve info from the device, apply configs
// to it, generate telemetry streams and program the RIB/FIB.
//
// The GetConfig service retrieves the configuration from a target device for
// the YANG paths specified.
//
//	output, err = xr.GetConfig(ctx, conn, yang, id)
//	if err != nil {
//		log.Fatalf("Could not get the config from %s, %v", targets.Routers[d].Host, err)
//	}
//
// ShowCmdJSONOutput and ShowCmdTextOutput services return show command outputs
// JSON encoded or as unstructured text correspondingly.
//
//	switch *enc {
//	case "text":
//		output, err = xr.ShowCmdTextOutput(ctx, conn, *cli, id)
//	case "json":
//		output, err = xr.ShowCmdJSONOutput(ctx, conn, *cli, id)
//	default:
//		log.Fatalf("Do NOT recognize encoding: %v\n", *enc)
//	}
//
// Tutorials to setup a testbed have been posted at:
//
// - Programming IOS-XR with gRPC and Go: https://xrdocs.github.io/programmability/tutorials/2017-08-04-programming-ios-xr-with-grpc-and-go/
//
// - Validate the intent of network config changes: https://xrdocs.github.io/programmability/tutorials/2017-08-14-validate-the-intent-of-network-config-changes/
//
package xrgrpc
