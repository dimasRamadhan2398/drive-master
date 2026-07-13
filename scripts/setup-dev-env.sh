#!/usr/bin/env bash
# scripts/setup-dev-env.sh
set -euo pipefail

SSH_TARGET="root@203.194.114.20"
TARGET_PATH="/var/www/drive"

echo "=== 1. Building and Syncing Configs & Binaries ==="
chmod 600 deploy_key

# Build backend binaries
./scripts/build-binaries.sh

# Sync files
./scripts/sync-binaries.sh

echo "=== 2. Setting up dev environment on remote server ==="
ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none "$SSH_TARGET" 'bash -s' << 'EOF'
set -euo pipefail

echo ">>> Creating dev.api subdomain in cPanel if it doesn't exist..."
if [ ! -d "/home/drivemaster/public_html/dev.api" ]; then
  /usr/bin/uapi --user=drivemaster SubDomain addsubdomain domain=dev.api rootdomain=drivemaster.id dir=public_html/dev.api
  echo "Subdomain dev.api.drivemaster.id created."
else
  echo "Subdomain directory already exists."
fi

# Ensure correct permissions for the subdomain directory
mkdir -p /home/drivemaster/public_html/dev.api
chown -R drivemaster:drivemaster /home/drivemaster/public_html/dev.api

echo ">>> Creating dev databases in PostgreSQL..."
DATABASES=(
  "dev_drivemaster_user_service"
  "dev_drivemaster_core_service"
  "dev_drivemaster_booking_service"
  "dev_drivemaster_payment_service"
)

for DB in "${DATABASES[@]}"; do
  echo "Checking database: $DB..."
  if sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw "$DB"; then
    echo "Database $DB already exists."
  else
    echo "Creating database $DB owned by admin_drive..."
    sudo -u postgres psql -c "CREATE DATABASE $DB OWNER admin_drive;"
  fi
done

echo ">>> Configuring dev Zookeeper & Kafka..."
# Create data directories
mkdir -p /var/lib/zookeeper/data-dev
mkdir -p /var/lib/kafka/data-dev
chown -R root:root /var/lib/zookeeper/data-dev
chown -R root:root /var/lib/kafka/data-dev

# Setup Zookeeper dev config
cp /opt/kafka/config/zookeeper.properties /opt/kafka/config/zookeeper-dev.properties
sed -i 's|clientPort=2181|clientPort=2182|g' /opt/kafka/config/zookeeper-dev.properties
sed -i 's|dataDir=/tmp/zookeeper|dataDir=/var/lib/zookeeper/data-dev|g' /opt/kafka/config/zookeeper-dev.properties
sed -i 's|dataDir=/var/lib/zookeeper/data|dataDir=/var/lib/zookeeper/data-dev|g' /opt/kafka/config/zookeeper-dev.properties

# Setup Kafka dev config
cp /opt/kafka/config/server.properties /opt/kafka/config/server-dev.properties
sed -i 's/broker.id=1/broker.id=2/g' /opt/kafka/config/server-dev.properties
sed -i 's/listeners=PLAINTEXT:\/\/0.0.0.0:9092/listeners=PLAINTEXT:\/\/0.0.0.0:9093/g' /opt/kafka/config/server-dev.properties
# Make sure listeners is active (not commented) and pointing to 9093
sed -i 's/#listeners=PLAINTEXT:\/\/:9092/listeners=PLAINTEXT:\/\/0.0.0.0:9093/g' /opt/kafka/config/server-dev.properties

if grep -q "advertised.listeners" /opt/kafka/config/server-dev.properties; then
  sed -i 's|advertised.listeners=.*|advertised.listeners=PLAINTEXT://localhost:9093|g' /opt/kafka/config/server-dev.properties
else
  echo "advertised.listeners=PLAINTEXT://localhost:9093" >> /opt/kafka/config/server-dev.properties
fi

sed -i 's/zookeeper.connect=localhost:2181/zookeeper.connect=localhost:2182/g' /opt/kafka/config/server-dev.properties
sed -i 's|log.dirs=/tmp/kafka-logs|log.dirs=/var/lib/kafka/data-dev|g' /opt/kafka/config/server-dev.properties
sed -i 's|log.dirs=/var/lib/kafka/data|log.dirs=/var/lib/kafka/data-dev|g' /opt/kafka/config/server-dev.properties

echo ">>> Registering & starting dev Zookeeper and Kafka in PM2..."
pm2 delete zookeeper-dev kafka-dev || true
KAFKA_HEAP_OPTS="-Xmx128m -Xms128m" pm2 start /opt/kafka/bin/zookeeper-server-start.sh \
  --name "zookeeper-dev" \
  -- /opt/kafka/config/zookeeper-dev.properties

echo "Waiting for zookeeper-dev to start..."
sleep 5

KAFKA_HEAP_OPTS="-Xmx256m -Xms256m" pm2 start /opt/kafka/bin/kafka-server-start.sh \
  --name "kafka-dev" \
  -- /opt/kafka/config/server-dev.properties

echo ">>> Starting dev microservices in PM2..."
cd /var/www/drive
pm2 delete dev-user-service dev-core-service dev-booking-service dev-payment-service dev-api-gateway || true
pm2 start ecosystem-dev.config.js

EOF

echo "=== 3. Uploading .htaccess for dev.api subdomain ==="
rsync -avz -e "ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none" ./scripts/htaccess_dev "$SSH_TARGET":/home/drivemaster/public_html/dev.api/.htaccess
ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none "$SSH_TARGET" "chown drivemaster:drivemaster /home/drivemaster/public_html/dev.api/.htaccess && chmod 644 /home/drivemaster/public_html/dev.api/.htaccess"

echo "=== Setup Dev Env successfully executed! ==="
