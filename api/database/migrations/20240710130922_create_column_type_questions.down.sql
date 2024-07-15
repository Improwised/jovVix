-- +migrate Down

-- remove column type
ALTER TABLE  questions
DROP COLUMN IF EXISTS type;