-- +migrate Up

-- rename user_session_questions to user_quiz_responses
ALTER TABLE IF EXISTS user_session_questions
RENAME TO user_quiz_responses;
