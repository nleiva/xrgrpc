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

- Clear text

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

- JSON

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

### CLI config

Apply cli commands to the device/router

```bash
examples/setconfig$ ./setconfig -cli "interface Lo11 ipv6 address 2001:db8::/128"
Config Applied
```

## XR Config

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

## Compiling the Go client example

You simply execute `go build`. Pre-compiled binaries for GOARCH="amd64" are also part of the repo.

- MAC OS X (GOOS=darwin): [xrgrpc](xrgrpc)
- Windows (GOOS=windows): [xrgrpc.exe](xrgrpc.exe)
