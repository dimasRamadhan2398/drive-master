#!/usr/bin/env bash
# scripts/setup-server-logging.sh
set -euo pipefail

SSH_TARGET="root@203.194.114.20"

echo "=== Running logging configuration on remote server ==="
ssh -i deploy_key -o StrictHostKeyChecking=no "$SSH_TARGET" 'bash -s' << 'EOF'
set -euo pipefail

echo "=== 1. Truncating current bloated PM2 logs ==="
truncate -s 0 /root/.pm2/logs/*.log || true
echo "Truncated existing PM2 logs."

echo "=== 2. Configuring PM2 Log Rotation ==="
# Install pm2-logrotate module
pm2 install pm2-logrotate

# Configure pm2-logrotate settings
pm2 set pm2-logrotate:max_size 10M
pm2 set pm2-logrotate:retain 5
pm2 set pm2-logrotate:compress true
pm2 set pm2-logrotate:workerInterval 30
pm2 set pm2-logrotate:rotateInterval '0 0 * * *'
echo "PM2 log rotation configured successfully."

echo "=== 3. Configuring Docker Log Rotation ==="
DOCKER_DAEMON_JSON='{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}'

mkdir -p /etc/docker
echo "$DOCKER_DAEMON_JSON" > /etc/docker/daemon.json
echo "Docker logging daemon.json written."

echo "=== 4. Restarting Docker service to apply changes ==="
systemctl restart docker
echo "Docker service restarted."

echo "=== Server Logging Setup Completed ==="
EOF
