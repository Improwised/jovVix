-- +migrate Up

-- create user_session and connect it with users
CREATE TABLE IF NOT EXISTS "user_sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" bpchar(20) NOT NULL,
  "is_host" bool NOT NULL DEFAULT false,
  "quiz_session_id" uuid NOT NULL,
  "leave_at" timestamp,
  "quiz_analysis" json,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

-- foreign key
ALTER TABLE "user_sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_sessions" ADD FOREIGN KEY ("quiz_session_id") REFERENCES "quiz_sessions" ("id");

-- index
CREATE INDEX quiz_session_id_user_id_idx
ON user_sessions (quiz_session_id);
