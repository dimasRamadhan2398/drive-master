-- Rollback: Drop member_profiles table

DROP INDEX IF EXISTS idx_member_profiles_user_id;
DROP TABLE IF EXISTS member_profiles;
