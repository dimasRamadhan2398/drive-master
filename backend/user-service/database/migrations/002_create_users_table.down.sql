-- Rollback: Drop users table

DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
