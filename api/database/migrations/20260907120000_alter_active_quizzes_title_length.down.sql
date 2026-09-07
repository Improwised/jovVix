-- +migrate Down

ALTER TABLE active_quizzes
ALTER COLUMN title TYPE varchar(30);
