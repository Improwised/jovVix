-- +migrate Down

-- remove not null
ALTER TABLE IF EXISTS questions
ALTER COLUMN duration_in_seconds DROP NOT NULL;
