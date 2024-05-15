-- +migrate Up

-- Rename column score to points
ALTER TABLE questions
RENAME COLUMN score TO points;
