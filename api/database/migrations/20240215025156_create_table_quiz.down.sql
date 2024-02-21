-- +migrate Down

-- remove fk
ALTER TABLE quizzes DROP CONSTRAINT quizzes_creator_id_fkey;

-- remove table
DROP TABLE IF EXISTS quizzes;
