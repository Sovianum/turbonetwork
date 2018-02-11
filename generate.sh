#!/usr/bin/env bash

PLUGIN=/home/artem/gowork/bin/protoc-gen-go

## generate node service
#protoc -I proto node_service.proto --plugin=${PLUGIN} --go_out=plugins=grpc:nodeservice/pb
## generate graph service
#protoc -I proto network_service.proto --plugin=${PLUGIN} --go_out=plugins=grpc:networkservice/pb

protoc -I proto node_service.proto network_service.proto --plugin=${PLUGIN} --go_out=plugins=grpc:pb