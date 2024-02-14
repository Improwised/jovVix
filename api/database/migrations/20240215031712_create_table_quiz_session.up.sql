-- +migrate Up

-- create quiz session and connect it with users and quiz
CREATE TABLE IF NOT EXISTS "quiz_sessions" (
  "id" uuid PRIMARY KEY,
  "code" integer NOT NULL,
  "title" varchar(30),
  "quiz_id" uuid NOT NULL,
  "admin_id" bpchar(20),
  "max_attempt" integer NOT NULL DEFAULT 1,
  "activated_to" timestamp,
  "activated_from" timestamp,
  "is_active" bool NOT NULL DEFAULT false,
  "quiz_analysis" json,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

-- foreign keys
ALTER TABLE "quiz_sessions" ADD FOREIGN KEY ("quiz_id") REFERENCES "quizzes" ("id");

ALTER TABLE "quiz_sessions" ADD FOREIGN KEY ("admin_id") REFERENCES "users" ("id");

-- indexes
CREATE INDEX session_code_is_active_idx
ON quiz_sessions (code, is_active);
