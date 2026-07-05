#!/bin/bash
set -e
echo "Setting up local development environment..."
command -v go     >/dev/null || { echo "Go required"; exit 1; }
command -v docker >/dev/null || { echo "Docker required"; exit 1; }
[ -f .env ] || cp .env.example .env
echo "Starting infrastructure..."
docker compose up -d postgres-user postgres-product postgres-order postgres-payment redis rabbitmq
echo "Waiting for databases to be ready..."
sleep 15
echo "Setup complete! Run 'make up' to start all services."
