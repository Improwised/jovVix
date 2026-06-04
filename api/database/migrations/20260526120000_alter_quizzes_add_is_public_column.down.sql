-- +migrate Down
ALTER TABLE quizzes
DROP COLUMN IF EXISTS is_public;
