-- +migrate Down

-- remove fk
ALTER TABLE quiz_questions
DROP CONSTRAINT IF EXISTS quiz_questions_next_question_fkey;

-- remove column
ALTER TABLE quiz_questions
DROP COLUMN next_question;