-- +migrate Up

-- make type column not null

ALTER TABLE questions
ALTER COLUMN "type" SET NOT NULL;