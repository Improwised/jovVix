-- +migrate Down

-- drop index
DROP INDEX IF EXISTS user_quiz_response_idx;

-- remove fk
ALTER TABLE user_quiz_responses
DROP CONSTRAINT IF EXISTS user_quiz_responses_user_played_quiz_id_fkey;

-- remove column user_played_quiz_id
ALTER TABLE user_quiz_responses
DROP COLUMN IF EXISTS user_played_quiz_id;

-- add column user_id
ALTER TABLE user_quiz_responses
ADD COLUMN user_id bpchar(20);
ALTER TABLE user_quiz_responses
ADD COLUMN active_quiz_id uuid;

-- add fk
ALTER TABLE user_quiz_responses
ADD CONSTRAINT user_session_questions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE user_quiz_responses
ADD CONSTRAINT user_session_questions_quiz_session_id_fkey FOREIGN KEY (active_quiz_id) REFERENCES active_quizzes (id);
