#!/usr/bin/env bash
# scripts/sync-databases.sh
set -euo pipefail

echo "=== Starting Database Sync (Local -> VPS) ==="

# SSH configuration
SSH_TARGET="root@203.194.114.20"
TARGET_DIR="/var/www/drive"

# 1. Create a local temp dumps directory inside the workspace
DUMP_DIR="./database_dumps"
mkdir -p "$DUMP_DIR"

DATABASES=(
  "drivemaster_user_service"
  "drivemaster_core_service"
  "drivemaster_booking_service"
  "drivemaster_payment_service"
)

# 2. Dump local PostgreSQL databases from the Docker container
for DB in "${DATABASES[@]}"; do
  echo "Dumping local database: $DB..."
  docker exec -i drive_master_postgres pg_dump -U admin_drive -d "$DB" --clean --no-owner --no-privileges > "$DUMP_DIR/$DB.sql"
done

# 3. Create target directory on VPS and upload SQL files
echo "Uploading dumps to VPS..."
ssh -i deploy_key -o StrictHostKeyChecking=no "$SSH_TARGET" "mkdir -p $TARGET_DIR/database_dumps"
rsync -avz -e "ssh -i deploy_key -o StrictHostKeyChecking=no" "$DUMP_DIR"/ "$SSH_TARGET":"$TARGET_DIR/database_dumps"/

# 4. Restore databases on VPS PostgreSQL
echo "Restoring databases on VPS..."
for DB in "${DATABASES[@]}"; do
  echo "Re-creating database $DB on VPS..."
  # Terminate active connections, drop and recreate
  ssh -i deploy_key -o StrictHostKeyChecking=no "$SSH_TARGET" "sudo -u postgres psql -c \"SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '$DB' AND pid <> pg_backend_pid();\""
  ssh -i deploy_key -o StrictHostKeyChecking=no "$SSH_TARGET" "sudo -u postgres psql -c \"DROP DATABASE IF EXISTS $DB;\""
  ssh -i deploy_key -o StrictHostKeyChecking=no "$SSH_TARGET" "sudo -u postgres psql -c \"CREATE DATABASE $DB OWNER admin_drive;\""
  
  echo "Importing schema and data for $DB..."
  ssh -i deploy_key -o StrictHostKeyChecking=no "$SSH_TARGET" "sudo -u postgres psql -d $DB -f $TARGET_DIR/database_dumps/$DB.sql"
done

# 5. Clean up temporary files
echo "Cleaning up temp files..."
rm -rf "$DUMP_DIR"
ssh -i deploy_key -o StrictHostKeyChecking=no "$SSH_TARGET" "rm -rf $TARGET_DIR/database_dumps"

echo "=== Database Sync Completed Successfully! ==="
