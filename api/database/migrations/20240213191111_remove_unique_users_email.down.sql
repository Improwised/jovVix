-- +migrate Down

-- add unique constraint again
ALTER TABLE users
ADD CONSTRAINT users_email_ukey UNIQUE (email);