#!/bin/sh
set -e

export CONFIG_PATH=/app/config.yaml

echo "Running database migrations..."
./user-service migrate

echo "Running database seeders..."
./user-service seed

echo "Starting User Service..."
./user-service serve