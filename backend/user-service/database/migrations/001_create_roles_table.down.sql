-- Rollback: Drop roles table

DROP INDEX IF EXISTS idx_roles_name;
DROP TABLE IF EXISTS roles;
