-- Rollback: Drop instructor_profiles table

DROP INDEX IF EXISTS idx_instructor_profiles_is_active;
DROP INDEX IF EXISTS idx_instructor_profiles_user_id;
DROP TABLE IF EXISTS instructor_profiles;
