#!/bin/sh

# Add longsleep/golang-backports PPA to your sources
sudo add-apt-repository ppa:longsleep/golang-backports -y
sudo apt-get update

# Install Go
sudo apt-get install golang-1.8-go -y

# Create a workspace directory.
mkdir /home/ubuntu/go
mkdir /home/ubuntu/go/src
mkdir /home/ubuntu/go/bin
mkdir /home/ubuntu/go/pkg

# Setup some environment variables and write them to .profile file.
echo "export GOROOT=/usr/lib/go-1.8/" >> /home/ubuntu/.profile
echo "export GOPATH=/home/ubuntu/go" >> /home/ubuntu/.profile
echo "export PATH=$PATH:/usr/lib/go-1.8/bin:/home/ubuntu/go/bin" >> /home/ubuntu/.profile
. /home/ubuntu/.profile

# Get the IOS XR gRPC library
go get github.com/nleiva/xrgrpc

# Just change ownership as previous commands were executed by root.
chown -R ubuntu:ubuntu /home/ubuntu/go/*