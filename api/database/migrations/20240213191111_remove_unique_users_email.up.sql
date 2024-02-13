-- +migrate Up
ALTER TABLE users
DROP CONSTRAINT users_username_ukey;
