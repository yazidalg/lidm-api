-- Add question_type column to questions table if it doesn't exist
ALTER TABLE questions ADD COLUMN question_type VARCHAR(20) DEFAULT 'regular' AFTER correct_answer;

-- Update existing questions to have regular type
UPDATE questions SET question_type = 'regular' WHERE question_type IS NULL OR question_type = '';
