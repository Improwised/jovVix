-- +migrate Up
CREATE TABLE IF NOT EXISTS shared_quizzes (
  id uuid PRIMARY KEY,
  quiz_id uuid references quizzes(id),
  shared_to VARCHAR(255) NOT NULL,
  shared_by bpchar(20) REFERENCES users(id),
  permission VARCHAR(50) DEFAULT 'read',
  created_at timestamp NOT NULL DEFAULT (now()),
  updated_at timestamp NOT NULL DEFAULT (now()),
  UNIQUE (quiz_id, shared_to)
);
