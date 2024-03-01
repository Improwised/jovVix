-- +migrate Up

ALTER TABLE quiz_sessions
ALTER COLUMN code DROP NOT NULL;
