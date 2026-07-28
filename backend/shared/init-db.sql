-- Role
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'admin_drive') THEN
        CREATE ROLE admin_drive WITH LOGIN SUPERUSER CREATEDB CREATEROLE PASSWORD 'drivemaster123';
    END IF;
END $$;

-- Databases (run outside any transaction)
CREATE DATABASE drivemaster_user_service         OWNER admin_drive;
CREATE DATABASE drivemaster_core_service         OWNER admin_drive;
CREATE DATABASE drivemaster_booking_service      OWNER admin_drive;
CREATE DATABASE drivemaster_catalog_service      OWNER admin_drive;
CREATE DATABASE drivemaster_payment_service      OWNER admin_drive;
CREATE DATABASE drivemaster_voucher_service      OWNER admin_drive;
CREATE DATABASE drivemaster_notification_service OWNER admin_drive;
SELECT 'Done!' AS status;