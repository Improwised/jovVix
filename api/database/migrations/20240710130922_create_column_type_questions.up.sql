-- +migrate Up

-- add column type
ALTER TABLE  questions
ADD COLUMN IF NOT EXISTS type SMALLINT;