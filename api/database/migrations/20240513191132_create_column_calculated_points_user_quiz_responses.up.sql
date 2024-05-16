-- +migrate Up

-- add column calculated_points
ALTER TABLE  user_quiz_responses
ADD COLUMN calculated_points int DEFAULT 0;
