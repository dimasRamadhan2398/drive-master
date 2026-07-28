-- Rollback: Drop instructor_areas table

DROP INDEX IF EXISTS idx_instructor_areas_area_id;
DROP INDEX IF EXISTS idx_instructor_areas_instructor_id;
DROP TABLE IF EXISTS instructor_areas;
