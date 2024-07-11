-- +migrate Up

-- add column type
ALTER TABLE  questions
ADD COLUMN type int4 NOT NULL;