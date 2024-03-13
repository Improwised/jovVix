-- +migrate Up

-- add column duration_in_seconds
ALTER TABLE questions
ADD COLUMN duration_in_seconds int;
