-- +migrate Up
ALTER TABLE users
ADD COLUMN username VARCHAR(12);

ALTER TABLE users
ADD CONSTRAINT users_username_ukey UNIQUE (username);
