-- +migrate Up

-- create table
CREATE TABLE IF NOT EXISTS "session_questions" (
  "id" uuid PRIMARY KEY,
  "question_id" uuid,
  "next_question" uuid,
  "quiz_session_id" uuid,
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
ADD FOREIGN KEY ("quiz_session_id")
REFERENCES "quiz_sessions" ("id");

-- create index
CREATE INDEX session_questions_quiz_session_idx
ON session_questions (quiz_session_id, question_id);
