-- +migrate Up
ALTER TABLE quizzes
ADD COLUMN category_id uuid REFERENCES quiz_categories (id) ON DELETE SET NULL,
ADD COLUMN cover_image TEXT;
