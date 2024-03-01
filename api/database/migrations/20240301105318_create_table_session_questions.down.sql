-- +migrate Down

-- remove index
DROP INDEX IF EXISTS session_questions_quiz_session_idx;

-- remove fk
ALTER TABLE "session_questions" DROP CONSTRAINT IF EXISTS session_questions_quiz_session_fkey;
ALTER TABLE "session_questions" DROP CONSTRAINT IF EXISTS session_questions_next_question_fkey;
ALTER TABLE "session_questions" DROP CONSTRAINT IF EXISTS session_questions_question_id_fkey;

-- remove table
DROP TABLE IF EXISTS "session_questions";
