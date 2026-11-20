-- +migrate Down
ALTER TABLE active_quizzes DROP COLUMN solo_completed;
ALTER TABLE active_quizzes DROP COLUMN quiz_version;
ALTER TABLE active_quizzes DROP COLUMN is_solo;
