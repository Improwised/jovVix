-- +migrate Up
CREATE TABLE quiz_rooms (
 id uuid PRIMARY KEY,
 admin_id bpchar(20) NOT NULL REFERENCES users(id),
 code integer NOT NULL UNIQUE CHECK (code BETWEEN 100000 AND 999999),
 current_session_id uuid NOT NULL REFERENCES active_quizzes(id),
 closed_at timestamp,
 expires_at timestamp NOT NULL DEFAULT now() + interval '24 hours',
 host_seen_at timestamp NOT NULL DEFAULT now(),
 created_at timestamp NOT NULL DEFAULT now()
);
ALTER TABLE active_quizzes ADD COLUMN room_id uuid REFERENCES quiz_rooms(id);
CREATE INDEX active_quizzes_room_idx ON active_quizzes(room_id);
CREATE TABLE quiz_room_members (
 room_id uuid NOT NULL REFERENCES quiz_rooms(id),
 user_id bpchar(20) NOT NULL REFERENCES users(id),
 player_name text NOT NULL,
 left_at timestamp,
 last_seen_at timestamp NOT NULL DEFAULT now(),
 PRIMARY KEY(room_id,user_id)
);
CREATE UNIQUE INDEX quiz_room_member_names ON quiz_room_members(room_id,lower(btrim(player_name)));
