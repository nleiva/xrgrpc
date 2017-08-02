[![GoDoc](https://godoc.org/github.com/nleiva/xrgrpc?status.svg)](https://godoc.org/github.com/nleiva/xrgrpc)

# gRPC library for Cisco IOS XR

Minimalistic library to interact with IOS XR devices using the gRPC framework. Look at the [IOS XR proto file](proto/ems_grpc.proto) for the description of the service interface and the structure of the payload messages. gRPC uses protocol buffers as the Interface Definition Language (IDL).

## Usage

CLI examples to use the library are provided in the [example](example/) folder. The CLI specified in the examples is not definitive and might change as we go.

### Get Config

Retrieves the config from one target device described in [config.json](example/input/config.json), for the YANG paths specified in [yangpaths.json](example/input/yangpaths.json)

```bash
example/getconfig$ ./getconfig

Config from [2001:420:2cff:1204::5502:1]:57344
{
 "data": {
  "Cisco-IOS-XR-ifmgr-cfg:interface-configurations": {
   "interface-configuration": [
    {
     "active": "act",
     "interface-name": "Loopback60",
     "interface-virtual": [
      null
     ],
     "Cisco-IOS-XR-ipv6-ma-cfg:ipv6-network": {
      "addresses": {
       "regular-addresses": {
        "regular-address": [
...
2017/07/21 15:11:47 This process took 1.195469855s
```

### Show Commands

Provides the output of IOS XR cli commands for one router defined in [config.json](example/input/config.json). Two output format options are available; Unstructured text and JSON encoded:

- **Clear text**

```bash
example/showcmd$ ./showcmd -cli "show isis database" -enc text

Output from [2001:420:2cff:1204::5502:1]:57344
 
----------------------------- show isis database ------------------------------

IS-IS BB2 (Level-2) Link State Database
LSPID                 LSP Seq Num  LSP Checksum  LSP Holdtime  ATT/P/OL
mrstn-5502-1.cisco.com.00-00* 0x0000000c   0x1558        3066            0/0/0
mrstn-5502-2.cisco.com.00-00  0x00000012   0x6e0c        3066            0/0/0
mrstn-5501-1.cisco.com.00-00  0x0000000c   0x65d5        1150            0/0/0

 Total Level-2 LSP count: 3     Local Level-2 LSP count: 1


2017/07/21 15:37:00 This process took 2.480039252s
```

- **JSON**

```bash
example/showcmd$ ./showcmd -cli "show isis database" -enc json

Config from [2001:420:2cff:1204::5502:1]:57344
 [{
 "Cisco-IOS-XR-clns-isis-oper:isis": {
<snip>
       {
        "system-id": "0151.0250.0002",
        "local-is-flag": false,
        "host-levels": "isis-levels-2",
        "host-name": "mrstn-5502-2.cisco.com"
       },
       {
        "system-id": "0151.0250.0003",
        "local-is-flag": false,
        "host-levels": "isis-levels-2",
        "host-name": "mrstn-5501-1.cisco.com"
       },
       {
        "system-id": "0151.0250.0001",
        "local-is-flag": true,
        "host-levels": "isis-levels-2",
        "host-name": "mrstn-5502-1.cisco.com"
...
2017/07/21 15:37:27 This process took 1.54038192s
```

### Configuring the router

- **CLI config** (Merge)

Applies CLI config commands on the device/router from the list in [config.json](example/input/config.json).

```bash
example/setconfig$ ./setconfig -cli "interface Lo11 ipv6 address 2001:db8::/128"

Config applied to [2001:420:2cff:1204::5502:1]:57344

2017/07/21 15:24:17 This process took 1.779449886s
```

You can verify the config on the router:

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#show run interface lo11
Fri Jul 21 15:24:24.199 EDT
interface Loopback11
 ipv6 address 2001:db8::/128
!
```

- **JSON** (Merge)

Applies YANG/JSON formatted config to one device/router (merges with existing config) from the list in [config.json](example/input/config.json). It reads the target from [yangconfig.json](example/input/yangconfig.json). The 

```bash
example/mergeconfig$ ./mergeconfig 

Config merged on [2001:420:2cff:1204::5502:1]:57344 -> Request ID: 8162, Response ID: 8162

2017/07/21 15:18:07 This process took 1.531427437s
```

You can verify the config on the router:

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#show run interface lo201
Fri Jul 21 15:18:24.046 EDT
interface Loopback201
 description New Loopback 201
 ipv6 address 2001:db8:20::1/128
!
```

- **JSON** (Replace)

Applies YANG/JSON formatted config to one device/router (replaces the config for this section) from the list in [config.json](example/input/config.json). It learns the config to replace from [yangconfigrep.json](example/input/yangconfigrep.json). If we had merged instead, we would have ended up with two IPv6 addresses in this example.

```bash
example/replaceconfig$ ./replaceconfig 

Config replaced on [2001:420:2cff:1204::5502:1]:57344 -> Request ID: 4616, Response ID: 4616

2017/07/21 15:21:27 This process took 1.623047025s
```

You can verify the config on the router:

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#show run interface lo201
Fri Jul 21 15:21:48.053 EDT
interface Loopback201
 description New Loopback 221
 ipv6 address 2001:db8:22::2/128
!
```

### Removing router config

- **JSON**

Removes YANG/JSON formatted config on one device/router from [config.json](example/input/config.json). It reads the config to delete from [yangdelconfig.json](example/input/yangdelconfig.json). The following example deletes both interfaces configured in the Merge example. See [yangdelintadd.json](example/input/yangdelintadd.json) to delete just the IP address and [yangdelintdesc.json](example/input/yangdelintdesc.json) for only the description of the interface.

```bash
example/deleteconfig$ ./deleteconfig 

Config Deleted on [2001:420:2cff:1204::5502:1]:57344 -> Request ID: 2856, Response ID: 2856

2017/07/21 15:06:46 This process took 730.329288ms
```

On the router:

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#show configuration commit changes 1000000039
Mon Jul 17 15:54:59.221 EDT
Building configuration...
!! IOS XR Configuration version = 6.2.2.22I
no interface Loopback201
no interface Loopback301
end
```

### **CLI config multiple routers simultaneously** (Merge)

Applies CLI config commands to the list of routers specified on [config.json](example/input/config.json). Notice that even though we added two devices, the execution time did NOT increase. This is possible because of the use of [Golang Concurrency](https://blog.golang.org/pipelines) primitives.

```bash
example/setconfiglist$ ./setconfiglist -cli "interface Lo33 ipv6 address 2001:db8:33::1/128"

Config applied to [2001:420:2cff:1204::5502:2]:57344



Config applied to [2001:420:2cff:1204::5501:1]:57344



Config applied to [2001:420:2cff:1204::5502:1]:57344


2017/07/21 15:32:11 This process took 1.773893901s
```

You can verify the config on the routers:

```
RP/0/RP0/CPU0:mrstn-5501-1.cisco.com#sh run int Lo33
Fri Jul 21 15:32:35.468 EDT
interface Loopback33
 ipv6 address 2001:db8:33::1/128
!
```

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#sh run int Lo33
Fri Jul 21 15:33:07.281 EDT
interface Loopback33
 ipv6 address 2001:db8:33::1/128
!
```

```
RP/0/RP0/CPU0:mrstn-5502-2.cisco.com#sh run int Lo33
Fri Jul 21 15:33:14.504 EDT
interface Loopback33
 ipv6 address 2001:db8:33::1/128
!
```


### Telemetry

- **JSON**

Subscribe to a Telemetry stream. The Telemetry message is defined in [telemetry.proto](proto/telemetry/telemetry.proto). The payload is JSON encoded (GPBKV).

```bash
example/telemetry$ ./telemetry -subs "LLDP"
Time 1500666991103, Path: Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail
{
  "NodeId": {
    "NodeIdStr": "mrstn-5502-1.cisco.com"
  },
  "Subscription": {
    "SubscriptionIdStr": "LLDP"
  },
  "encoding_path": "Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail",
  "collection_id": 1,
  "collection_start_time": 1500666991103,
  "msg_timestamp": 1500666991103,
  "data_gpbkv": [
    {
      "timestamp": 1500666991108,
      "ValueByType": null,
      "fields": [
...
```

The Subscription ID has to exist on the device <sup>[1](#myfootnote1)</sup>.

```
telemetry model-driven
 sensor-group LLDPNeighbor
  sensor-path Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail
 !
 subscription LLDP
  sensor-group-id LLDPNeighbor sample-interval 15000
 !
!
```

- **GPB**

Again, we subscribe to a Telemetry stream but we request the content is encoded with [protobuf](https://developers.google.com/protocol-buffers/). To decode the message we need to look at the "LLDP neighbor details" definition in [lldp_neighbor.proto](proto/telemetry/lldp/lldp_neighbor.proto). We parse the message and modify the output to illustrate how to access to each field on it.

```bash
example/telemetrygpb$ ./telemetrygpb -subs "LLDP"
Time 1500667512299, Path: Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail
{
  "node_name": "0/RP0/CPU0",
  "interface_name": "HundredGigE0/0/0/22",
  "device_id": "mrstn-5502-2.cisco.com"
}
Type:  6.2.2.22I, NCS-5500, Address value:"2001:558:2::2"  

{
  "node_name": "0/RP0/CPU0",
  "interface_name": "HundredGigE0/0/0/21",
  "device_id": "mrstn-5502-2.cisco.com"
}
Type:  6.2.2.22I, NCS-5500, Address value:"2001:558:2::2"  

{
  "node_name": "0/RP0/CPU0",
  "interface_name": "HundredGigE0/0/0/1",
  "device_id": "mrstn-5502-2.cisco.com"
}
Type:  6.2.2.22I, NCS-5500, Address value:"2001:f00:bb::2"
...
```

The Subscription ID has to exist on the device <sup>[1](#myfootnote1)</sup>.

```
telemetry model-driven
 sensor-group LLDPNeighbor
  sensor-path Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail
 !
 subscription LLDP
  sensor-group-id LLDPNeighbor sample-interval 15000
 !
!
```
<a name="myfootnote1">[1]</a>: [gNMI](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi.proto) defines a variant where you do not need this config.

### Telemetry

You can manually define the target without the config file [config.json](example/input/config.json), by calling the functional options "WithValue". See the snippet below from [definetarget](example/definetarget/main.go).

```go
// Manually specify target parameters.
router, err := xr.BuildRouter(
	xr.WithUsername("cisco"),
	xr.WithPassword("cisco"),
	xr.WithHost("[2001:420:2cff:1204::5502:2]:57344"),
	xr.WithCreds("../input/ems5502-2.pem"),
	xr.WithTimeout(5),
)
```

## XR gRPC Config

The following is the configuration required on the IOS XR device in order to enable gRPC dial-in with TLS support.

```
!! IOS XR Configuration version = 6.2.2
grpc
 port 57344
 tls
 !
 address-family ipv6
!
```

### Port range

While you can select any not-used port on the device, it's recommended to choose one from the 57344-57999 range.

```
mrstn-5502-1 emsd: [1058]: %MGBL-EMS-4-EMSD_PORT_RANGE : The configured port 56500 is outside of the range of [57344, 57999]. It will consume an additional LPTS entry.
```

## Certificate file

You need to retrive the `ems.pem` file from the IOS XR device (after enabling gRPC/TLS) and put it in the [input](example/input) folder (or any other location specified in [config.json](example/input/config.json)). You can find the file in the router on either `/misc/config/grpc/` or `/var/xr/config/grpc`.

```bash
[xrrouter.cisco.com:/var/xr/config/grpc]$ ls -la
total 20
drwxr-xr-x  3 root root 4096 Jul  5 17:47 .
drwxr-xr-x 10 root root 4096 Jul  3 12:50 ..
drwx------  2 root root 4096 Jul  3 12:50 dialout
-rw-------  1 root root 1675 Jul  5 17:47 ems.key
-rw-rw-rw-  1 root root 1513 Jul  5 17:47 ems.pem
[xrrouter.cisco.com:/var/xr/config/grpc]$
```

## Compiling the proto files

The Go generated code in [ems_grpc.pb.go](proto/ems/ems_grpc.pb.go) is the result of the following:

```bash
proto/ems/$ protoc --go_out=plugins=grpc:. ems_grpc.proto
```

The Go generated code in [lldp_neighbor.pb.go](proto/telemetry/lldp/lldp_neighbor.pb.go) is the result of the following:

```bash
proto/telemetry/lldp$ protoc --go_out=. lldp_neighbor.proto 
```

## Compiling the Examples

Simply execute `go build` on the corresponding example folder. E.g.

```bash
example/telemetry$ go build
```