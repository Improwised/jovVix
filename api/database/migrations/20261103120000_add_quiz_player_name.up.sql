-- +migrate Up
ALTER TABLE user_played_quizzes ADD COLUMN player_name VARCHAR(12);
CREATE UNIQUE INDEX user_played_quizzes_session_player_name_unique
ON user_played_quizzes (active_quiz_id, lower(btrim(player_name)))
WHERE player_name IS NOT NULL AND is_host = false;
