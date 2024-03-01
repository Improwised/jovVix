-- +migrate Up

-- rename quiz_sessions to active_quizzes
ALTER TABLE quiz_sessions
RENAME TO active_quizzes;

-- rename columns
ALTER TABLE user_sessions
RENAME COLUMN quiz_session_id TO active_quiz_id;

ALTER TABLE user_session_questions
RENAME COLUMN quiz_session_id TO active_quiz_id;
