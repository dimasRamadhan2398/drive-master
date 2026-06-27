#!/usr/bin/env bash
# scripts/build-binaries.sh
set -euo pipefail

echo "=== Starting Local Compilation ==="

SERVICES=("user-service" "core-service" "booking-service" "payment-service")

for SERVICE in "${SERVICES[@]}"; do
  echo "Building $SERVICE..."
  cd "backend/$SERVICE"
  
  # Ensure the bin directory exists
  mkdir -p bin
  
  # Compile for Linux amd64 (CGO disabled for static linking)
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor -ldflags="-w -s" -o bin/"$SERVICE" ./cmd
  
  echo "$SERVICE build completed."
  cd ../..
done

echo "=== All microservices compiled successfully! ==="
