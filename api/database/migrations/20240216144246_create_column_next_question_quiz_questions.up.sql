-- +migrate Up

-- add next question column
ALTER TABLE quiz_questions
ADD COLUMN next_question uuid;

-- add foreign key
ALTER TABLE quiz_questions
ADD CONSTRAINT quiz_questions_next_question_fkey
FOREIGN KEY (next_question)
REFERENCES questions (id);
