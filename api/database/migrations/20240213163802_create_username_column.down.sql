-- +migrate Down

-- drop index
DROP INDEX users_username_hash_idx;

-- drop constraint
ALTER TABLE users
DROP CONSTRAINT users_username_ukey;

-- drop column
ALTER TABLE users
DROP COLUMN username;

