-- +migrate Down
ALTER TABLE user_quiz_responses
DROP COLUMN IF EXISTS option_order;
