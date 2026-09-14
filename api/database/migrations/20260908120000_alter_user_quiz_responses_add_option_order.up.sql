-- +migrate Up
-- Maps the option position shown to a participant to the original question option key.
ALTER TABLE user_quiz_responses
ADD COLUMN option_order json;
