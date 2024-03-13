-- +migrate Up

-- rename column is_count to is_attend
ALTER TABLE user_quiz_responses 
RENAME COLUMN is_count TO is_attend;

-- change default to false
ALTER TABLE user_quiz_responses 
ALTER COLUMN is_attend SET DEFAULT false;

