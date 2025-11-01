-- Add simple position tracking fields to users table
-- Note: Run each ALTER statement separately if some columns already exist

-- Add position_type column
ALTER TABLE users ADD COLUMN position_type VARCHAR(20) DEFAULT 'stable';

-- Add previous_position column  
ALTER TABLE users ADD COLUMN previous_position INT DEFAULT 0;