-- Migration: Create work_experiences table
-- Description: Stores work experience records for instructors

CREATE TABLE IF NOT EXISTS work_experiences (
    id SERIAL PRIMARY KEY,
    instructor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    company_name VARCHAR(100) NOT NULL,
    role VARCHAR(100) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    description TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_work_experiences_instructor_id ON work_experiences(instructor_id);
CREATE INDEX IF NOT EXISTS idx_work_experiences_company_name ON work_experiences(company_name);
