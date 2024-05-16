-- +migrate Down

-- remove column calculated_points
ALTER TABLE  user_quiz_responses
DROP COLUMN calculated_points;
