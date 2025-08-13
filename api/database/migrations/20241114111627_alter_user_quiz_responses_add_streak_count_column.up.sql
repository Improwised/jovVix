-- +migrate Up
ALTER TABLE user_quiz_responses
ADD COLUMN streak_count INTEGER DEFAULT 0;
