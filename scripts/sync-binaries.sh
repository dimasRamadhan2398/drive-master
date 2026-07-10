#!/usr/bin/env bash
# scripts/sync-binaries.sh
set -euo pipefail

echo "=== Syncing Binaries & Assets to VPS ==="

# SSH Connection config
SSH_TARGET="root@203.194.114.20"
TARGET_PATH="/var/www/drive"

chmod 600 deploy_key

SSH_OPTS="-i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none -o ControlMaster=auto -o ControlPath=./ssh-control-%r@%h:%p -o ControlPersist=10m"

# 1. Sync root configs
echo "Syncing root configs..."
rsync -avz -e "ssh $SSH_OPTS" ./docker-compose-vps.yml "$SSH_TARGET":"$TARGET_PATH"/docker-compose.yml
rsync -avz -e "ssh $SSH_OPTS" ./package.json ./package-lock.json ./.env.production ./ecosystem.config.js ./ecosystem-dev.config.js "$SSH_TARGET":"$TARGET_PATH"/

# 2. Sync frontend build output (.output/)
echo "Syncing frontend built output..."
ssh $SSH_OPTS "$SSH_TARGET" "mkdir -p $TARGET_PATH/frontend"
rsync -avz --delete -e "ssh $SSH_OPTS" ./.output/ "$SSH_TARGET":"$TARGET_PATH"/frontend/.output/

# 3. Sync backend services (including Go binaries, config files, and migrations; excluding Go source code/vendor to speed up transfer)
echo "Syncing backend services..."
rsync -avz \
  --exclude="node_modules" \
  --exclude=".git" \
  --exclude=".claude" \
  --exclude="tmp" \
  --exclude="vendor/" \
  --exclude="*.go" \
  --exclude="go.mod" \
  --exclude="go.sum" \
  --exclude="Dockerfile" \
  --exclude="docker-compose.yml" \
  --exclude="entrypoint.sh" \
  -e "ssh $SSH_OPTS" ./backend/ "$SSH_TARGET":"$TARGET_PATH"/backend/

# Close control socket
ssh $SSH_OPTS -O exit "$SSH_TARGET" 2>/dev/null || true
rm -f ./ssh-control-* 2>/dev/null || true

echo "=== Syncing completed! ==="
