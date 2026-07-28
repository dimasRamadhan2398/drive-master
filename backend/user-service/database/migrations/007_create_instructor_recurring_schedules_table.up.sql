-- Create instructor_recurring_schedules table
-- Stores recurring time slots for instructors (e.g., every Mon-Fri 09:00-10:00, 13:00-14:00)
CREATE TABLE instructor_recurring_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6), -- 0=Sunday, 1=Monday, ..., 6=Saturday
    start_time VARCHAR(10) NOT NULL, -- Format: HH:MM (e.g., "09:00")
    end_time VARCHAR(10) NOT NULL,   -- Format: HH:MM (e.g., "10:00")
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Each instructor can only have one slot per time slot per day
    CONSTRAINT unique_instructor_day_time UNIQUE (instructor_id, day_of_week, start_time, end_time)
);

-- Index for fast lookups by instructor
CREATE INDEX idx_recurring_schedules_instructor ON instructor_recurring_schedules(instructor_id);

-- Index for day_of_week filtering
CREATE INDEX idx_recurring_schedules_day ON instructor_recurring_schedules(day_of_week);
