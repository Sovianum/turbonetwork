#!/usr/bin/env bash

PLUGIN=/home/artem/gowork/bin/protoc-gen-go

# generate node service
protoc -I nodeservice/pb service.proto --plugin=${PLUGIN} --go_out=plugins=grpc:nodeservice/pb