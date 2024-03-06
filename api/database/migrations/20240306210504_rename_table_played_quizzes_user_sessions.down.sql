-- +migrate Down

-- rename user_played_quizzes to user_sessions
ALTER TABLE IF EXISTS user_played_quizzes
RENAME TO user_sessions;
