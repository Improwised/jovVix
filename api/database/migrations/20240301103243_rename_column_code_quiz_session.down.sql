-- +migrate Down

-- rename code column
ALTER TABLE quiz_sessions
RENAME COLUMN invitation_code TO code;
