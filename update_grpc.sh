#!/bin/bash
set -xe

# ensure old generated files are removed
find "$(dirname ${0})" -type f -name '*.pb.go' -exec rm -f {} +

# update grpc files
protoc -I shepard --go_out=plugins=grpc:shepard shepard/*.proto
