-- +migrate Up

-- add column question_delivery_time
ALTER TABLE active_quizzes
ADD COLUMN question_delivery_time timestamp;
