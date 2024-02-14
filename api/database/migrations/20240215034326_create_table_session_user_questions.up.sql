-- +migrate Up

-- create user_session_questions table and give keys
CREATE TABLE "user_session_questions" (
  "id" uuid PRIMARY KEY,
  "user_id" bpchar(20) NOT NULL,
  "quiz_session_id" uuid NOT NULL,
  "question_id" uuid NOT NULL,
  "answers" json,
  "calculated_score" int DEFAULT 0,
  "is_count" bool DEFAULT true,
  "response_time" integer DEFAULT -1,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);


-- foreign keys
ALTER TABLE "user_session_questions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_session_questions" ADD FOREIGN KEY ("quiz_session_id") REFERENCES "quiz_sessions" ("id");

ALTER TABLE "user_session_questions" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id");

-- index

CREATE INDEX user_session_questions_idx
ON user_session_questions (quiz_session_id);
