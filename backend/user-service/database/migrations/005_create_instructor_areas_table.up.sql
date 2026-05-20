-- Migration: Create instructor_areas table
-- Description: Stores coverage areas for instructors (references external Area service)

CREATE TABLE IF NOT EXISTS instructor_areas (
    instructor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    area_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (instructor_id, area_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_instructor_areas_instructor_id ON instructor_areas(instructor_id);
CREATE INDEX IF NOT EXISTS idx_instructor_areas_area_id ON instructor_areas(area_id);
