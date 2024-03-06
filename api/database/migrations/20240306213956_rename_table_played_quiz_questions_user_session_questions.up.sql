-- +migrate Up

-- rename user_session_questions to played_quiz_questions
ALTER TABLE IF EXISTS user_session_questions
RENAME TO user_quiz_response;
