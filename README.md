# gRPC library for Cisco IOS XR

Minimalistic library to interact with IOS XR devices using the gRPC framework. Look at the [IOS XR proto file](proto/ems_grpc.proto) for the description of the service interface and the structure of the payload messages. gRPC uses protocol buffers as the Interface Definition Language (IDL).

## Usage

CLI examples to use the library are provided in the [examples](examples/) folder. The CLI specified is not definitive and will most likely change as we go.

### Get Config

Retrieves the config from the target device described in [config.json](examples/input/config.json) for the YANG paths specified in [yangpaths.json](examples/input/yangpaths.json)

```bash
examples/getconfig$ ./getconfig
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
```

### Show Commands

Provides the output of IOS XR cli commands on the router defined in [config.json](examples/input/config.json). Two output format options are available; Unstructured text and JSON encoded:

- **Clear text**

```bash
examples/showcmd$ ./showcmd -cli "show isis database" -enc text

----------------------------- show isis database ------------------------------

IS-IS BB2 (Level-2) Link State Database
LSPID                 LSP Seq Num  LSP Checksum  LSP Holdtime  ATT/P/OL
mrstn-5502-1.cisco.com.00-00* 0x0000000b   0x9c44        1395            0/0/0
mrstn-5502-2.cisco.com.00-00  0x0000000c   0x863f        1564            0/0/0

 Total Level-2 LSP count: 2     Local Level-2 LSP count: 1

$
```

- **JSON**

```bash
examples/showcmd$ ./showcmd -cli "show isis database" -enc json
[{
 "Cisco-IOS-XR-clns-isis-oper:isis": {
<snip>
{
 "Cisco-IOS-XR-clns-isis-oper:isis": {
  "instances": {
   "instance": [
    {
     "instance-name": "BB2",
     "host-names": {
      "host-name": [
       {
        "system-id": "0151.0250.0002",
        "local-is-flag": false,
        "host-levels": "isis-levels-2",
        "host-name": "mrstn-5502-2.cisco.com"
       },
       {
        "system-id": "0151.0250.0001",
        "local-is-flag": true,
        "host-levels": "isis-levels-2",
        "host-name": "mrstn-5502-1.cisco.com"
       }
      ]
     }
...
$
```

### Configuring the router

- **CLI config** (Merge)

Applies CLI config commands on the device/router.

```bash
examples/setconfig$ ./setconfig -cli "interface Lo11 ipv6 address 2001:db8::/128"
Config Applied
```

On the router:

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#show run interface lo11
Mon Jul 17 11:33:28.065 EDT
interface Loopback11
 ipv6 address 2001:db8::/128
!
```

- **JSON** (Merge)

Applies YANG/JSON formatted config to the device/router (merges with existing config). It reads the target from [yangconfig.json](examples/input/yangconfig.json). The 

```bash
examples/mergeconfig$ ./mergeconfig 
Config Applied -> Request ID: 163, Response ID: 163
```

On the router:

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#show run interface lo201
Mon Jul 17 15:06:22.521 EDT
interface Loopback201
 description New Loopback 201
 ipv6 address 2001:db8:20::1/128
!
```

- **JSON** (Replace)

Applies YANG/JSON formatted config to the device/router (replaces the config for this section). It reads the info from [yangconfigrep.json](examples/input/yangconfigrep.json). If we had merged instead, we would have ended up with two IPv6 addresses in this example.

```bash
examples/replaceconfig$ ./replaceconfig 
Config Replaced -> Request ID: 543, Response ID: 543
```

On the router:

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#show run int lo201
Mon Jul 17 17:06:13.376 EDT
interface Loopback201
 description New Loopback 221
 ipv6 address 2001:db8:22::2/128
!
```

### Removing router config

- **JSON**

Removes YANG/JSON formatted config on the device/router. It reads the config to delete from [yangdelconfig.json](examples/input/yangdelconfig.json). The follwowing example deletes both interfaces configured in the Merge example. See [yangdelintadd.json](examples/input/yangdelintadd.json) to delete just the IP address and [yangdelintdesc.json](examples/input/yangdelintdesc.json) for only the description of the interface.

```bash
examples/deleteconfig$ ./deleteconfig 
Config Deleted -> Request ID: 236, Response ID: 236
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

### Telemetry

Subscribe to a Telemetry stream. The Telemetry message is defined in [telemetry.proto](proto/telemetry/telemetry.proto). The payload is JSON encoded, we will add an example encoded with GPB.

```bash
examples/telemetry$ ./telemetry -subs "LLDP"
Time 1500576676957, Path: Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail
{
  "NodeId": {
    "NodeIdStr": "mrstn-5502-1.cisco.com"
  },
  "Subscription": {
    "SubscriptionIdStr": "LLDP"
  },
  "encoding_path": "Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail",
  "collection_id": 117,
  "collection_start_time": 1500576676957,
  "msg_timestamp": 1500576676957,
  "collection_end_time": 1500576676968
}
...
```

The Subscription ID has to exist on the device.

```
telemetry model-driven
 sensor-group LLDP
  sensor-path Cisco-IOS-XR-ethernet-lldp-oper:lldp/nodes/node/neighbors/details/detail
 !
```

## XR gRPC Config

The following is the configuration requiered on the IOS XR device in order to enable gRPC dial-in with TLS support.

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

While you can select any not-used port on the device, it's recomended to chose one from the 57344-57999 range.

```
mrstn-5502-1 emsd: [1058]: %MGBL-EMS-4-EMSD_PORT_RANGE : The configured port 56500 is outside of the range of [57344, 57999]. It will consume an additional LPTS entry.
```

## Certificate file

You need to retrive the `ems.pem` file from the IOS XR device (after enabling gRPC/TLS) and put it in the [input](examples/input) folder (or any other location specified in [config.json](examples/input/config.json)). You can find the file in the router on either `/misc/config/grpc/` or `/var/xr/config/grpc`.

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

## Compiling the proto file

The Go generated code in [ems_grpc.pb.go](proto/ems_grpc.pb.go) is the result of the following:

```bash
$ protoc --go_out=plugins=grpc:. ems_grpc.proto
```

## Compiling the Examples

Simply execute `go build` on each folder.
