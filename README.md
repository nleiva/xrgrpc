# gRPC library for Cisco IOS XR

Minimalistic library to interact with IOS XR devices using the gRPC framework. Look at the [IOS XR proto file](proto/ems_grpc.proto) for message and service definitions.

## Usage

A CLI example is provided to use the library. This is not definitive, it will change as we go.

### Get Config

Retrieves the config from the target device defined in [config.json](examples/input/config.json) for the YANG paths specified in [yangpaths.json](examples/input/yangpaths.json)

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

Provides the output for IOS XR cli commands on the router defined in [config.json](examples/input/config.json). Two encoding options available:

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

Apply cli commands to the device/router

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

Applies YANG/JSON formatted config to the device/router. It reads the info from [yangconfig.json](examples/input/yangconfig.json).

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

### Removing router config

- **JSON**

Removes YANG/JSON formatted config on the device/router. It reads the config to delete from [yangdelconfig.json](examples/input/yangdelconfig.json).

```bash
examples/deleteconfig$ ./deleteconfig 
Config Deleted -> Request ID: 236, Response ID: 236
nleiva@~/go/src/github.com/nleiva/xrgrpc/examples/deleteconfig$ 
```

On the router:

```
RP/0/RP0/CPU0:mrstn-5502-1.cisco.com#show configuration commit changes 1000000039
Mon Jul 17 15:54:59.221 EDT
Building configuration...
!! IOS XR Configuration version = 6.2.2.22I-Lindt
no interface Loopback201
no interface Loopback301
end
```


## XR gRPC Config

The following is the configuration requiered on the IOS XR device in order to enable gRPC dial-in.

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
