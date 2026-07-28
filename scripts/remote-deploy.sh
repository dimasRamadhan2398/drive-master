#!/usr/bin/env bash
set -euo pipefail

# remote-deploy.sh
# Arguments taken from environment: SSH_TARGET_PATH
TARGET=${SSH_TARGET_PATH:-/var/www/drive}
cd "$TARGET"

echo "Starting remote deploy in $TARGET"

# Ensure docker-compose is present
if ! command -v docker-compose >/dev/null 2>&1 && ! docker compose version >/dev/null 2>&1; then
  echo "docker-compose not found. Please install Docker Compose on the server."
  exit 1
fi

# Copy environment files if present (you should manage them securely outside this script)
if [ -f .env ]; then
  echo "Found .env in target"
fi

# Bring up services
if docker compose version >/dev/null 2>&1; then
  docker compose down || true
  docker compose pull || true
  docker compose up -d --build
else
  docker-compose down || true
  docker-compose pull || true
  docker-compose up -d --build
fi

# For Nuxt frontend, ensure build dir exists (if prebuilt was rsynced)
if [ -d frontend/.output ]; then
  echo "Frontend .output found"
else
  echo "No frontend build found. Building on server (this may take a while)"
  if [ -f package.json ]; then
    npm ci
    npm run build
  else
    echo "No package.json; skipping frontend build"
  fi
fi

# Post-deploy checks (basic)
if docker ps --filter "name=api-gateway" --format "{{.Names}}" | grep -q api-gateway; then
  echo "api-gateway is running"
else
  echo "Warning: api-gateway may not be running"
fi

echo "Deploy finished"
