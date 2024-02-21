-- +migrate Down

-- remove index
DROP INDEX user_session_questions_idx;

-- remove fk
ALTER TABLE user_session_questions DROP CONSTRAINT user_session_questions_question_id_fkey;
ALTER TABLE user_session_questions DROP CONSTRAINT user_session_questions_quiz_session_id_fkey;
ALTER TABLE user_session_questions DROP CONSTRAINT user_session_questions_user_id_fkey;

-- remove table
DROP TABLE IF EXISTS user_session_questions;
