package models

import (
	"database/sql"
	"testing"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

// GitHub #38: the reports list shows a session of my quiz that someone I shared it with hosted.
func TestListQuizzesAnalysisIncludesSharedQuizSessions(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`CREATE TABLE quizzes(id text PRIMARY KEY, creator_id text, title text, description text);
CREATE TABLE active_quizzes(id text PRIMARY KEY, quiz_id text, admin_id text, activated_from timestamp, activated_to timestamp);
CREATE TABLE active_quiz_questions(active_quiz_id text, question_id text);
CREATE TABLE user_played_quizzes(id text PRIMARY KEY, active_quiz_id text, is_host boolean);
CREATE TABLE user_quiz_responses(user_played_quiz_id text, question_id text, calculated_score int);
INSERT INTO quizzes VALUES
  ('00000000-0000-4000-8000-0000000000a1', 'owner', 'Owner quiz', NULL),
  ('00000000-0000-4000-8000-0000000000a2', 'other', 'Other quiz', NULL);
INSERT INTO active_quizzes VALUES
  ('00000000-0000-4000-8000-000000000001', '00000000-0000-4000-8000-0000000000a1', 'owner', '2026-11-20 10:00:00', '2026-11-20 10:30:00'),
  ('00000000-0000-4000-8000-000000000002', '00000000-0000-4000-8000-0000000000a1', 'colleague', '2026-11-20 11:00:00', '2026-11-20 11:30:00'),
  ('00000000-0000-4000-8000-000000000003', '00000000-0000-4000-8000-0000000000a2', 'colleague', '2026-11-20 12:00:00', '2026-11-20 12:30:00');
INSERT INTO active_quiz_questions VALUES
  ('00000000-0000-4000-8000-000000000001', 'q1'),
  ('00000000-0000-4000-8000-000000000002', 'q1'),
  ('00000000-0000-4000-8000-000000000003', 'q1');
INSERT INTO user_played_quizzes VALUES ('p1', '00000000-0000-4000-8000-000000000001', 0), ('p2', '00000000-0000-4000-8000-000000000002', 0), ('p3', '00000000-0000-4000-8000-000000000003', 0);
INSERT INTO user_quiz_responses VALUES ('p1', 'q1', 1), ('p2', 'q1', 0), ('p3', 'q1', 1);`); err != nil {
		t.Fatal(err)
	}
	model := InitQuizModel(goqu.New("sqlite3", db))

	sessions := func(user string) map[string]bool {
		t.Helper()
		rows, count, err := model.ListQuizzesAnalysis("", "", "", "", user, 1)
		if err != nil {
			t.Fatal(err)
		}
		if int(count) != len(rows) {
			t.Fatalf("%s: count %d, rows %d", user, count, len(rows))
		}
		got := map[string]bool{}
		for _, r := range rows {
			got[r.ID.String()] = true
		}
		return got
	}

	owner := sessions("owner")
	if len(owner) != 2 || !owner["00000000-0000-4000-8000-000000000001"] || !owner["00000000-0000-4000-8000-000000000002"] {
		t.Fatalf("owner sees %v, want both sessions of the owner's quiz", owner)
	}
	colleague := sessions("colleague")
	if len(colleague) != 2 || colleague["00000000-0000-4000-8000-000000000001"] {
		t.Fatalf("colleague sees %v, want only the two sessions they hosted", colleague)
	}
}
