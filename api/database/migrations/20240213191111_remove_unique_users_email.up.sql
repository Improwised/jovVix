-- +migrate Up

-- drop email unique constraint as currently user can directly enter in quiz
ALTER TABLE users
DROP CONSTRAINT users_email_key;
