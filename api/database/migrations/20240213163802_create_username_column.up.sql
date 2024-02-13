-- +migrate Up

-- add column with not null
ALTER TABLE users
ADD COLUMN username VARCHAR(12) Not null;

-- add unique constraint
ALTER TABLE users
ADD CONSTRAINT users_username_ukey UNIQUE (username);

-- add hash index
CREATE INDEX users_username_hash_idx 
ON users USING hash (username);
