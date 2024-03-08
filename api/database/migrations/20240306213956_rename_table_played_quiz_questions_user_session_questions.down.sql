-- +migrate Down

-- rename user_quiz_responses to user_session_questions
ALTER TABLE IF EXISTS user_quiz_responses
RENAME TO user_session_questions;
