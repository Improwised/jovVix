-- +migrate Down

-- remove column duration_in_seconds
ALTER TABLE IF EXISTS quiz_questions
DROP COLUMN IF EXISTS duration_in_seconds;
