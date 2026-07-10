#!/usr/bin/env bash
# scripts/migrate-docker-to-pm2.sh
set -euo pipefail

SSH_TARGET="root@203.194.114.20"

echo "=== Starting Migration from Docker to PM2 on Remote Server ==="
ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none "$SSH_TARGET" 'bash -s' << 'EOF'
set -euo pipefail

echo "=== Step 1: Installing Java OpenJDK 17 ==="
dnf install -y java-17-openjdk-headless
java -version

echo "=== Step 2: Downloading Apache Kafka ==="
mkdir -p /opt
if [ ! -d "/opt/kafka" ]; then
  echo "Downloading Kafka binary archive..."
  wget -c -q https://archive.apache.org/dist/kafka/3.6.1/kafka_2.13-3.6.1.tgz -O /opt/kafka.tgz
  echo "Extracting Kafka..."
  tar -xzf /opt/kafka.tgz -C /opt
  mv /opt/kafka_2.13-3.6.1 /opt/kafka
  rm -f /opt/kafka.tgz
  echo "Kafka installed to /opt/kafka."
else
  echo "Kafka is already installed in /opt/kafka."
fi

echo "=== Step 3: Configuring Kafka and Zookeeper ==="
# Reduce log retention hours to 72 hours to save disk space
sed -i 's/log.retention.hours=168/log.retention.hours=72/g' /opt/kafka/config/server.properties

echo "=== Step 4: Registering Zookeeper & Kafka in PM2 ==="
# Stop and delete existing pm2 instances of the same name to prevent duplicates
pm2 delete zookeeper kafka || true

# Start Zookeeper under PM2 (limiting JVM to 128M heap)
KAFKA_HEAP_OPTS="-Xmx128m -Xms128m" pm2 start /opt/kafka/bin/zookeeper-server-start.sh \
  --name "zookeeper" \
  -- /opt/kafka/config/zookeeper.properties

# Wait a few seconds for Zookeeper to start
sleep 5

# Start Kafka under PM2 (limiting JVM to 256M heap)
KAFKA_HEAP_OPTS="-Xmx256m -Xms256m" pm2 start /opt/kafka/bin/kafka-server-start.sh \
  --name "kafka" \
  -- /opt/kafka/config/server.properties

echo "PM2 Services started successfully."
pm2 list

EOF
