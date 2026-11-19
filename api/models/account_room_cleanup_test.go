package models

import (
	"database/sql"
	"testing"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

// Exercise real foreign keys, including the room/session cycle, without touching
// application accounts or requiring an external database.
func TestAccountDeletionRoomCleanup(t *testing.T) {
	for _, target := range []string{"host", "player", "both"} {
		for _, rollback := range []bool{false, true} {
			name := target + "/commit"
			if rollback {
				name = target + "/rollback"
			}
			t.Run(name, func(t *testing.T) {
				db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
				if err != nil {
					t.Fatal(err)
				}
				defer db.Close()
				db.SetMaxOpenConns(1)
				schema := `CREATE TABLE users(id text PRIMARY KEY);
CREATE TABLE active_quizzes(id text PRIMARY KEY,admin_id text REFERENCES users(id),room_id text REFERENCES quiz_rooms(id));
CREATE TABLE quiz_rooms(id text PRIMARY KEY,admin_id text REFERENCES users(id),current_session_id text REFERENCES active_quizzes(id));
CREATE TABLE quiz_room_members(room_id text REFERENCES quiz_rooms(id),user_id text REFERENCES users(id));
CREATE TABLE user_played_quizzes(id text PRIMARY KEY,active_quiz_id text REFERENCES active_quizzes(id));
CREATE TABLE user_quiz_responses(user_played_quiz_id text REFERENCES user_played_quizzes(id));
CREATE TABLE active_quiz_questions(active_quiz_id text REFERENCES active_quizzes(id));
INSERT INTO users VALUES('host'),('player'),('both'),('other');
INSERT INTO active_quizzes(id,admin_id) VALUES('h','host'),('b','both'),('o','other');
INSERT INTO quiz_rooms VALUES('rh','host','h'),('rb','both','b'),('ro','other','o');
UPDATE active_quizzes SET room_id='r'||id;
INSERT INTO quiz_room_members VALUES('rh','player'),('rh','both'),('rb','player'),('ro','both'),('ro','other');
INSERT INTO user_played_quizzes VALUES('ph','h'),('pb','b'),('po','o');
INSERT INTO user_quiz_responses VALUES('ph'),('pb'),('po');
INSERT INTO active_quiz_questions VALUES('h'),('b'),('o');`
				if _, err = db.Exec(schema); err != nil {
					t.Fatal(err)
				}
				database := goqu.New("sqlite3", db)
				tx, err := database.Begin()
				if err != nil {
					t.Fatal(err)
				}
				defer tx.Rollback()
				model := InitActiveQuizModel(database, nil)
				if err = model.DeleteActiveQuizzesAndRelatedDataByUserId(tx, target); err != nil {
					t.Fatal(err)
				}
				if _, err = tx.Delete("users").Where(goqu.Ex{"id": target}).Executor().Exec(); err != nil {
					t.Fatalf("account deletion blocked by foreign key: %v", err)
				}
				if rollback {
					err = tx.Rollback()
				} else {
					err = tx.Commit()
				}
				if err != nil {
					t.Fatal(err)
				}
				var rooms, users int
				if err = db.QueryRow("SELECT count(*) FROM quiz_rooms").Scan(&rooms); err != nil {
					t.Fatal(err)
				}
				if err = db.QueryRow("SELECT count(*) FROM users").Scan(&users); err != nil {
					t.Fatal(err)
				}
				wantRooms, wantUsers := 3, 4
				if !rollback {
					wantUsers = 3
					if target != "player" {
						wantRooms = 2
					}
				}
				if rooms != wantRooms || users != wantUsers {
					t.Fatalf("rooms/users=%d/%d want %d/%d", rooms, users, wantRooms, wantUsers)
				}
				var other int
				if err = db.QueryRow("SELECT count(*) FROM quiz_room_members WHERE room_id='ro' AND user_id='other'").Scan(&other); err != nil || other != 1 {
					t.Fatalf("unrelated membership changed: %d %v", other, err)
				}
			})
		}
	}
}
