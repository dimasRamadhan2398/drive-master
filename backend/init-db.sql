-- Role
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'admin_drive') THEN
        CREATE ROLE admin_drive WITH LOGIN SUPERUSER CREATEDB CREATEROLE PASSWORD 'drivemaster123';
    END IF;
END $$;

-- Databases (run outside any transaction)
CREATE DATABASE drivemaster_user_service;
CREATE DATABASE drivemaster_core_service;
CREATE DATABASE drivemaster_booking_service;
CREATE DATABASE drivemaster_catalog_service;
CREATE DATABASE drivemaster_payment_service;
CREATE DATABASE drivemaster_voucher_service;
CREATE DATABASE drivemaster_notification_service;


GRANT ALL PRIVILEGES ON DATABASE drivemaster_user_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_core_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_booking_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_catalog_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_payment_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_voucher_service TO admin_drive;
GRANT ALL PRIVILEGES ON DATABASE drivemaster_notification_service TO admin_drive;

-- Print completion message
SELECT 'Drive Master Database Initialization Complete!' AS status;