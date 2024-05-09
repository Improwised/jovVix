-- +migrate Down

-- Rename column if exists
ALTER TABLE IF EXISTS user_quiz_responses
RENAME COLUMN is_attend TO is_count;

-- Change default to true
ALTER TABLE IF EXISTS user_quiz_responses
ALTER COLUMN is_count SET DEFAULT true;
