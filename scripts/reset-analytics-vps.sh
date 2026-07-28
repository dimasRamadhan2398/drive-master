#!/usr/bin/env bash
# scripts/reset-analytics-vps.sh
set -euo pipefail

SSH_TARGET="root@203.194.114.20"

echo "=== Resetting Analytics and Transactional Data on VPS ==="

# Run psql commands on each database
ssh -i deploy_key -o StrictHostKeyChecking=no -o IPQoS=none "$SSH_TARGET" 'bash -s' << 'EOF'
set -euo pipefail

# List of databases to clean
DBS=(
  "drivemaster_payment_service"
  "dev_drivemaster_payment_service"
  "drivemaster_booking_service"
  "dev_drivemaster_booking_service"
  "drivemaster_core_service"
  "dev_drivemaster_core_service"
  "drivemaster_user_service"
  "dev_drivemaster_user_service"
)

# SQL block to safely truncate tables if they exist
SQL_RESET=$(cat << 'SQL'
DO $$
DECLARE
    t text;
BEGIN
    FOR t IN 
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
          AND table_name IN (
            'refunds', 'transactions', 'payments', 
            'driving_sessions', 'user_certifications', 'schedules', 
            'transaction_items', 'enrollments', 'user_entitlements', 
            'sales', 'sale_items', 'monthly_sales', 
            'entitlements', 'certifications'
          )
    LOOP
        RAISE NOTICE 'Truncating table: %', t;
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(t) || ' CASCADE';
    END LOOP;
END $$;
SQL
)

echo ">>> Truncating transactional/analytics tables in PostgreSQL databases..."
for db in "${DBS[@]}"; do
  echo "Cleaning database: $db..."
  sudo -u postgres psql -d "$db" -c "$SQL_RESET"
done

echo ">>> Flushing Redis cache..."
redis-cli flushall

echo ">>> Restarting all PM2 services..."
pm2 restart all

echo ">>> PM2 Status:"
pm2 list

EOF

echo "=== Analytics Reset Completed Successfully! ==="
