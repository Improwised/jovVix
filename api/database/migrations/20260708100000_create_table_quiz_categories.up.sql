-- +migrate Up

CREATE TABLE IF NOT EXISTS "quiz_categories" (
  "id" uuid PRIMARY KEY,
  "name" varchar(50) NOT NULL UNIQUE,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);
