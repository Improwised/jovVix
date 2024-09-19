-- +migrate Up
ALTER TABLE questions
ADD COLUMN question_media VARCHAR(10) DEFAULT 'text';

ALTER TABLE questions
ADD COLUMN options_media VARCHAR(10) DEFAULT 'text';

ALTER TABLE questions
ADD COLUMN resource TEXT DEFAULT '';