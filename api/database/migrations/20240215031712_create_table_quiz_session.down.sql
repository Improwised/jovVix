-- +migrate Down

-- remove index
DROP INDEX session_code_is_active_idx;

-- remove fk
ALTER TABLE quiz_sessions DROP CONSTRAINT quiz_sessions_admin_id_fkey;
ALTER TABLE quiz_sessions DROP CONSTRAINT quiz_sessions_quiz_id_fkey;

-- remove table
DROP TABLE IF EXISTS quiz_sessions;
