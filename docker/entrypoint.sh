#!/usr/bin/env bash

set -e
set -m # enable job control

go run ./cmd/start-rpc-servers &

echo "Waiting for RPC servers to start..."
sleep 5

echo "Creating channels..."
go run ./cmd/create-channels

fg %1 # Bring the RPC servers to the foreground to keep the container running
