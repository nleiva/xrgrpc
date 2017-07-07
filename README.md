# Testing gRPC
gRPC library for Cisco IOS XR

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