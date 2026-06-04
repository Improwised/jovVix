-- +migrate Up
ALTER TABLE quizzes
ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT false;
