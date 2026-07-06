#!/bin/bash
set -e
echo "Setting up local dev environment..."
command -v go     >/dev/null 2>&1 || { echo "ERROR: Go is not installed"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "ERROR: Docker is not installed"; exit 1; }
[ -f .env ] || cp .env.example .env
echo "Starting infrastructure..."
docker compose up -d postgres-user postgres-product postgres-order postgres-payment redis rabbitmq
echo "Waiting 15s for databases..."
sleep 15
echo "Done! Run: docker compose up -d"
