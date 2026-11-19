-- +migrate Down
DROP TABLE quiz_room_members;
ALTER TABLE active_quizzes DROP COLUMN room_id;
DROP TABLE quiz_rooms;
