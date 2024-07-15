-- +migrate Down

-- remove column type

UPDATE questions
SET "type" = NULL
WHERE "type" IS NOT NULL;