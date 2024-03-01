-- +migrate Up

ALTER TABLE quiz_sessions
DROP COLUMN max_attempt;
