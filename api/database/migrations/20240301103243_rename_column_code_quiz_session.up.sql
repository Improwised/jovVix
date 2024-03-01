-- +migrate Up

-- rename code column
ALTER TABLE quiz_sessions
RENAME COLUMN code TO invitation_code;
