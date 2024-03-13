-- +migrate Down

-- add column question_delivery_time
ALTER TABLE IF EXISTS active_quizzes
DROP COLUMN IF EXISTS question_delivery_time;
