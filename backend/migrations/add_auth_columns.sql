-- Add password and role columns to users table
-- This migration adds authentication support to the users table

-- Add password column (hashed password, never store plain text)
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255) NOT NULL DEFAULT '';

-- Add role column (default 'user', can be 'user' or 'admin')
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'user';

-- Create index on email for faster login queries
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Update existing users to have a default password (should be changed on first login)
-- In production, you would require users to reset their passwords
UPDATE users SET password = '$2a$10$defaulthashedpassword' WHERE password = '';
