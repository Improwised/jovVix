-- +migrate Down

-- rename active_quiz_questions to session_questions
ALTER TABLE IF EXISTS active_quiz_questions
RENAME TO session_questions;
