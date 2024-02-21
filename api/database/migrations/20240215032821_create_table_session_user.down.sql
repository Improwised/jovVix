-- +migrate Down

-- remove index
DROP INDEX quiz_session_id_user_id_idx;

-- remove fk
ALTER TABLE user_sessions DROP CONSTRAINT user_sessions_quiz_session_id_fkey;
ALTER TABLE user_sessions DROP CONSTRAINT user_sessions_user_id_fkey;

-- remove table
DROP TABLE IF EXISTS user_sessions;
