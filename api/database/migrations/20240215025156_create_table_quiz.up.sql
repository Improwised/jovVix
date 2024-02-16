-- +migrate Up

-- create quizzes table and add foreign-key
CREATE TABLE IF NOT EXISTS "quizzes" (
  "id" uuid PRIMARY KEY,
  "title" varchar(50) NOT NULL,
  "description" varchar(150),
  "creator_id" bpchar(20),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "quizzes" ADD FOREIGN KEY ("creator_id") REFERENCES "users" ("id");
