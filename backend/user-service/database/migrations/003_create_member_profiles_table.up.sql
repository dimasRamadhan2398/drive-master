-- Migration: Create member_profiles table
-- Description: Stores member-specific profile information

CREATE TABLE IF NOT EXISTS member_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    sessions_completed INTEGER DEFAULT 0,
    training_time INTEGER DEFAULT 0,
    average_rating DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for user lookups
CREATE INDEX IF NOT EXISTS idx_member_profiles_user_id ON member_profiles(user_id);
