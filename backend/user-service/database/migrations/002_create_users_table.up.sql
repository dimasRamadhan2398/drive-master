-- Migration: Create users table
-- Description: Stores user account information

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email_address VARCHAR(150) NOT NULL UNIQUE,
    phone_number VARCHAR(20),
    date_of_birth TIMESTAMP,
    image VARCHAR(255),
    address VARCHAR(500),
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email_address);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
