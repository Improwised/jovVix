-- +migrate Down
ALTER TABLE users
ADD CONSTRAINT users_email_ukey UNIQUE (email);