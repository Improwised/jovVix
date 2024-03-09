-- +migrate Up

-- rename session_questions to active_quiz_questions
ALTER TABLE IF EXISTS session_questions
RENAME TO active_quiz_questions;
