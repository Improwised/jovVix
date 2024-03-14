-- +migrate Down

-- remove not null
ALTER TABLE questions
ALTER COLUMN duration_in_seconds DROP NOT NULL;