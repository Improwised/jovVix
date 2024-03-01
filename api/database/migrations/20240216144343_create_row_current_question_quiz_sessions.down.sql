-- +migrate Down

-- remove fk
ALTER TABLE quiz_sessions
DROP CONSTRAINT IF EXISTS quiz_sessions_current_question_fkey;
