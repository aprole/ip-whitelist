#!/bin/sh

mkdir -p pb
protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    -I ./proto \
    --go_opt=Mproto/service.proto=github.com/aprole/ip-whitelist/pb \
    ./proto/service.proto
