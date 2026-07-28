-- Rollback: Drop work_experiences table

DROP INDEX IF EXISTS idx_work_experiences_company_name;
DROP INDEX IF EXISTS idx_work_experiences_instructor_id;
DROP TABLE IF EXISTS work_experiences;
