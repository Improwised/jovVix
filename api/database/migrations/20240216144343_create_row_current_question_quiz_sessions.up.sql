-- +migrate Up

-- add column current_question and is_question_active
ALTER TABLE quiz_sessions
ADD COLUMN current_question uuid;
ALTER TABLE quiz_sessions
ADD COLUMN is_question_active bool default false;

-- add foreign key
ALTER TABLE quiz_sessions
ADD CONSTRAINT quiz_sessions_current_question_fkey
FOREIGN KEY (current_question)
REFERENCES questions (id);
