-- +migrate Up

-- create question table
CREATE TABLE IF NOT EXISTS "questions" (
  "id" uuid PRIMARY KEY,
  "question" text NOT NULL,
  "options" json NOT NULL,
  "answers" json NOT NULL,
  "score" integer,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

-- connect question table with quiz
CREATE TABLE IF NOT EXISTS "quiz_questions" (
  "id" uuid PRIMARY KEY,
  "question_id" uuid,
  "quiz_id" uuid,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);


-- foreign keys
ALTER TABLE "quiz_questions" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id");
ALTER TABLE "quiz_questions" ADD FOREIGN KEY ("quiz_id") REFERENCES "quizzes" ("id");

-- index
CREATE INDEX quiz_questions_idx
ON quiz_questions (quiz_id, question_id);
