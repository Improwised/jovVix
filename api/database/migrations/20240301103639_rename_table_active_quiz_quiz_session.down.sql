-- +migrate Down

-- rename active_quizzes to quiz_sessions
ALTER TABLE active_quizzes
RENAME TO quiz_sessions;

-- rename columns
ALTER TABLE user_sessions
RENAME COLUMN active_quiz_id TO quiz_session_id;

ALTER TABLE user_session_questions
RENAME COLUMN active_quiz_id TO quiz_session_id;
