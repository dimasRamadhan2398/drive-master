#!/usr/bin/env bash
# scripts/build-binaries.sh
set -euo pipefail

echo "=== Starting Local Compilation ==="

SERVICES=("user-service" "core-service" "booking-service" "payment-service" "api-gateway")

for SERVICE in "${SERVICES[@]}"; do
  echo "Building $SERVICE..."
  cd "backend/$SERVICE"
  
  # Ensure the bin directory exists
  mkdir -p bin

  # Use -mod=vendor if vendor/ exists, otherwise use -mod=mod
  if [ -d "vendor" ]; then
    MOD_FLAG="-mod=vendor"
  else
    echo "  No vendor/ found for $SERVICE, downloading dependencies..."
    go mod download
    MOD_FLAG="-mod=mod"
  fi
  
  # Compile for Linux amd64 (CGO disabled for static linking)
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $MOD_FLAG -ldflags="-w -s" -o bin/"$SERVICE" ./cmd
  
  echo "$SERVICE build completed."
  cd ../..
done

echo "=== All microservices compiled successfully! ==="
