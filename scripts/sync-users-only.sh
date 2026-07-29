#!/usr/bin/env bash
# scripts/sync-users-only.sh
set -euo pipefail

echo "=== Starting User Data Sync (Local -> VPS) ==="

# SSH configuration
SSH_TARGET="root@203.194.114.20"
TARGET_DIR="/var/www/drive"

# 1. Create a local temp dumps directory inside the workspace
DUMP_DIR="./database_dumps"
mkdir -p "$DUMP_DIR"

DB="drivemaster_user_service"

echo "Dumping local user data tables..."
# Dump only users, roles, member_profiles, and instructor_profiles tables
docker exec -i drive_master_postgres pg_dump -U admin_drive -d "$DB" \
  -t public.users \
  -t public.roles \
  -t public.member_profiles \
  -t public.instructor_profiles \
  --clean --no-owner --no-privileges > "$DUMP_DIR/users_only.sql"

# 2. Create target directory on VPS and upload SQL file
echo "Uploading dump to VPS..."
ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none "$SSH_TARGET" "mkdir -p $TARGET_DIR/database_dumps"
rsync -avz -e "ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none" "$DUMP_DIR"/ "$SSH_TARGET":"$TARGET_DIR/database_dumps"/

# 3. Restore tables on VPS PostgreSQL
echo "Restoring user data tables on VPS..."
ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none "$SSH_TARGET" "sudo -u postgres psql -d $DB -f $TARGET_DIR/database_dumps/users_only.sql"

# 4. Clean up temporary files
echo "Cleaning up temp files..."
rm -rf "$DUMP_DIR"
ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none "$SSH_TARGET" "rm -rf $TARGET_DIR/database_dumps"

echo "=== User Data Sync Completed Successfully! ==="
