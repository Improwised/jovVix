-- +migrate Down
ALTER TABLE quizzes
DROP COLUMN IF EXISTS category_id,
DROP COLUMN IF EXISTS cover_image;
