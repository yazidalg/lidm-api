-- Migration untuk menambahkan field streak ke users table

-- Menambahkan kolom streak ke users table
ALTER TABLE users ADD COLUMN current_streak INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN max_streak INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN last_active_date TIMESTAMP;

-- Update existing users dengan default values
UPDATE users SET current_streak = 0, max_streak = 0 WHERE current_streak IS NULL;

-- Index untuk performa query streak
CREATE INDEX idx_users_current_streak ON users(current_streak);
CREATE INDEX idx_users_last_active_date ON users(last_active_date);
