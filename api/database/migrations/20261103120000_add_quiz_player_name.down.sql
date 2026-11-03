-- +migrate Down
DROP INDEX IF EXISTS user_played_quizzes_session_player_name_unique;
ALTER TABLE user_played_quizzes DROP COLUMN player_name;
