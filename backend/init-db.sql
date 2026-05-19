-- =====================================================
-- Drive Master Database Initialization Script
-- PostgreSQL 16+
-- =====================================================

-- Create admin role if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'admin_drive') THEN
        CREATE ROLE admin_drive WITH
            LOGIN SUPERUSER CREATEDB CREATEROLE
            PASSWORD 'drivemaster123';
        RAISE NOTICE 'Role admin_drive created successfully';
    ELSE
        RAISE NOTICE 'Role admin_drive already exists, skipping';
    END IF;
END
$$;

-- Grant privileges to admin role
GRANT ALL PRIVILEGES ON SCHEMA public TO admin_drive;
GRANT ALL PRIVILEGES ON SCHEMA extensions TO admin_drive;

-- Create databases for each microservice
DO $$
DECLARE
    dbs TEXT[] := ARRAY[
        'drivemaster_user_service',
        'drivemaster_core_service',
        'drivemaster_booking_service',
        'drivemaster_catalog_service',
        'drivemaster_payment_service',
        'drivemaster_voucher_service',
        'drivemaster_notification_service'
    ];
    db TEXT;
BEGIN
    FOREACH db IN ARRAY dbs LOOP
        IF NOT EXISTS (SELECT FROM pg_database WHERE datname = db) THEN
            EXECUTE format('CREATE DATABASE %I OWNER admin_drive', db);
            RAISE NOTICE 'Created database: %', db;
        ELSE
            RAISE NOTICE 'Database already exists, skipping: %', db;
        END IF;
    END LOOP;
END
$$;

-- Grant privileges on all databases to admin role
GRANT ALL PRIVILEGES ON DATABASE drivemaster_user_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_core_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_booking_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_catalog_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_payment_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_voucher_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_notification_service TO admin_drive;

-- Print completion message
SELECT 'Drive Master Database Initialization Complete!' AS status;