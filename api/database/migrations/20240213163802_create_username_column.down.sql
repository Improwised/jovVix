-- +migrate Down
ALTER TABLE users
DROP CONSTRAINT users_username_ukey;

ALTER TABLE users
DROP COLUMN username;

