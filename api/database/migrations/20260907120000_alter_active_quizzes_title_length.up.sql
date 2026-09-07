-- +migrate Up

-- A quiz title is copied into active_quizzes when a session starts. Keep this
-- capacity aligned with quizzes.title so valid quizzes can always be started.
ALTER TABLE active_quizzes
ALTER COLUMN title TYPE varchar(50);
