-- +migrate Up
ALTER TABLE active_quizzes ADD COLUMN is_solo boolean NOT NULL DEFAULT false;
ALTER TABLE active_quizzes ADD COLUMN quiz_version text NOT NULL DEFAULT '';
ALTER TABLE active_quizzes ADD COLUMN solo_completed boolean NOT NULL DEFAULT false;
