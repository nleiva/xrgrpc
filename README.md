# Testing gRPC
gRPC library for Cisco IOS XR

## Usage

This is not definitive, will change as we go.

```bash
$ go run grpccli.go -e text -c "show isis database"

----------------------------- show isis database ------------------------------

IS-IS BB2 (Level-2) Link State Database
LSPID                 LSP Seq Num  LSP Checksum  LSP Holdtime  ATT/P/OL
mrstn-5502-1.cisco.com.00-00* 0x0000000b   0x9c44        1395            0/0/0
mrstn-5502-2.cisco.com.00-00  0x0000000c   0x863f        1564            0/0/0

 Total Level-2 LSP count: 2     Local Level-2 LSP count: 1

$
```

```bash
$ go run grpccli.go -e json -c "show isis database"
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
    }
   ]
  }
 }
}
]
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
 port 56500
 tls
 !
 address-family ipv6
!
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