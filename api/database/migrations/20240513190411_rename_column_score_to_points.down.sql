-- +migrate Down

-- Rename column if exists
ALTER TABLE IF EXISTS questions
RENAME COLUMN IF EXISTS points TO score;
