-- +migrate Down

-- alter type column to allow null

ALTER TABLE questions
ALTER COLUMN "type" DROP NOT NULL;
