-- +migrate Down

-- add not null to email
ALTER TABLE users
ALTER COLUMN email SET NOT NULL;
