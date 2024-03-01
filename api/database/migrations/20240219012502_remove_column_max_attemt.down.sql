-- +migrate Down

-- create column
ALTER TABLE quiz_sessions
ADD COLUMN max_attempt int;