-- +migrate Down

-- rename played_quiz_questions to user_session_questions
ALTER TABLE IF EXISTS user_quiz_response
RENAME TO user_session_questions;
