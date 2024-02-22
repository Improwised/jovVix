-- +migrate Down

-- remove index
DROP INDEX quiz_questions_idx;

-- remove fk
ALTER TABLE quiz_questions DROP CONSTRAINT quiz_questions_quiz_id_fkey;
ALTER TABLE quiz_questions DROP CONSTRAINT quiz_questions_question_id_fkey;

-- remove table
DROP TABLE IF EXISTS quiz_questions;
DROP TABLE IF EXISTS questions;
