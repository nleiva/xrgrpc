# gRPC library for Cisco IOS XR

Minimalistic library to interact with IOS XR devices using the gRPC framework. Look at the [IOS XR proto file](proto/ems_grpc.proto) with message and service definitions.

## Usage

This is not definitive, will change as we go. Router config parameters are defined in [config.json](config.json).

- Clear text

```bash
$ ./xrgrpc -cli "show isis database" -enc text

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
$ ./xrgrpc -cli "show isis database" -enc json
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

## Compiling the proto file

```bash
$protoc --go_out=plugins=grpc:. ems_grpc.proto
```

## XR Config

```
!! IOS XR Configuration version = 6.2.2
grpc
 port 57344
 tls
 !
 address-family ipv6
!
```

Port range

```
mrstn-5502-1 emsd: [1058]: %MGBL-EMS-4-EMSD_PORT_RANGE : The configured port 56500 is outside of the range of [57344, 57999]. It will consume an additional LPTS entry.
```

## Cert file

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