#!/bin/bash
set -e
echo "Generating gRPC code from proto files..."
for svc in user product order; do
  protoc --go_out=. --go_opt=paths=source_relative \
         --go-grpc_out=. --go-grpc_opt=paths=source_relative \
         proto/$svc/$svc.proto
  echo "Generated $svc proto"
done
echo "Done!"
