-- Migration: Create instructor_profiles table
-- Description: Stores instructor-specific profile information

CREATE TABLE IF NOT EXISTS instructor_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    license_number VARCHAR(50) NOT NULL,
    license_expiry TIMESTAMP,
    bnsp_certificate_number VARCHAR(50) NOT NULL,
    bio TEXT,
    photo_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    number_of_students INTEGER DEFAULT 0,
    years_of_experience INTEGER DEFAULT 0,
    sessions_completed INTEGER DEFAULT 0,
    average_rating DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for user lookups
CREATE INDEX IF NOT EXISTS idx_instructor_profiles_user_id ON instructor_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_instructor_profiles_is_active ON instructor_profiles(is_active);
