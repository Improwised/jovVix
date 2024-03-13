-- +migrate Up

-- rename user_sessions to user_played_quizzes
ALTER TABLE IF EXISTS user_sessions
RENAME TO user_played_quizzes;
