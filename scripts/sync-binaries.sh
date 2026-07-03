#!/usr/bin/env bash
# scripts/sync-binaries.sh
set -euo pipefail

echo "=== Syncing Binaries & Assets to VPS ==="

# SSH Connection config
SSH_TARGET="root@203.194.114.20"
TARGET_PATH="/var/www/drive"

chmod 600 deploy_key

# 1. Sync root configs
echo "Syncing root configs..."
rsync -avz -e "ssh -i deploy_key -o StrictHostKeyChecking=no" ./docker-compose-vps.yml "$SSH_TARGET":"$TARGET_PATH"/docker-compose.yml
rsync -avz -e "ssh -i deploy_key -o StrictHostKeyChecking=no" ./package.json ./package-lock.json ./.env.production ./ecosystem.config.js "$SSH_TARGET":"$TARGET_PATH"/

# 2. Sync frontend build output (.output/)
echo "Syncing frontend built output..."
ssh -i deploy_key -o StrictHostKeyChecking=no "$SSH_TARGET" "mkdir -p $TARGET_PATH/frontend"
rsync -avz --delete -e "ssh -i deploy_key -o StrictHostKeyChecking=no" ./.output/ "$SSH_TARGET":"$TARGET_PATH"/frontend/.output/

# 3. Sync backend services (including Go binaries in bin/ and configs/refs)
echo "Syncing backend services..."
rsync -avz --exclude="node_modules" --exclude=".git" --exclude=".claude" --exclude="tmp" -e "ssh -i deploy_key -o StrictHostKeyChecking=no" ./backend/ "$SSH_TARGET":"$TARGET_PATH"/backend/

echo "=== Syncing completed! ==="
