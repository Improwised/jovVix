-- +migrate Up

-- drop index
DROP INDEX IF EXISTS user_session_questions_idx;

-- remove fk
ALTER TABLE user_quiz_responses
DROP CONSTRAINT user_session_questions_user_id_fkey;
ALTER TABLE user_quiz_responses
DROP CONSTRAINT user_session_questions_quiz_session_id_fkey;

-- remove column user_id
ALTER TABLE user_quiz_responses
DROP COLUMN user_id;
ALTER TABLE user_quiz_responses
DROP COLUMN active_quiz_id;

-- add column
ALTER TABLE user_quiz_responses
ADD COLUMN user_played_quiz_id uuid NOT NULL;

-- add fk
ALTER TABLE user_quiz_responses
ADD CONSTRAINT user_quiz_responses_user_played_quiz_id_fkey
FOREIGN KEY (user_played_quiz_id)
REFERENCES user_played_quizzes(id);

-- add index
CREATE INDEX user_quiz_response_idx
ON user_quiz_responses (user_played_quiz_id);
