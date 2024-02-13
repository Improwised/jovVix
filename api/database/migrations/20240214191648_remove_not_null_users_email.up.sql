-- +migrate Up

-- remove not null from email
ALTER TABLE users
ALTER COLUMN email DROP NOT NULL;
