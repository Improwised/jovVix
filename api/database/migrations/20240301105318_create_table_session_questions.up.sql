-- +migrate Up

-- create table
CREATE TABLE IF NOT EXISTS "session_questions" (
  "id" uuid PRIMARY KEY,
  "question_id" uuid,
  "next_question" uuid,
  "active_quiz_id" uuid,
  "order_no" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

-- create fk
ALTER TABLE "session_questions"
ADD FOREIGN KEY ("question_id")
REFERENCES "questions" ("id");

ALTER TABLE "session_questions"
ADD FOREIGN KEY ("next_question")
REFERENCES "questions" ("id");

ALTER TABLE "session_questions"
ADD FOREIGN KEY ("active_quiz_id")
REFERENCES "active_quizzes" ("id");

-- create index
CREATE INDEX session_questions_active_quizzes_idx
ON session_questions (active_quiz_id, question_id);
