package v1

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func TestSoloHostOptionOrder(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	session, participant, question := uuid.New(), uuid.New(), uuid.New()
	for _, statement := range []string{
		`CREATE TABLE active_quizzes(id text, admin_id text)`,
		`CREATE TABLE user_played_quizzes(id text, active_quiz_id text, user_id text,is_host boolean)`,
		`CREATE TABLE user_quiz_responses(id text,user_played_quiz_id text,question_id text,option_order text)`,
		`CREATE TABLE active_quiz_questions(active_quiz_id text,question_id text)`,
		`CREATE TABLE questions(id text, options text,type int)`,
	} {
		if _, err = db.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}
	db.Exec(`INSERT INTO active_quizzes VALUES(?, 'host')`, session.String())
	db.Exec(`INSERT INTO user_played_quizzes VALUES(?,?, 'host',false)`, participant.String(), session.String())
	db.Exec(`INSERT INTO active_quiz_questions VALUES(?,?)`, session.String(), question.String())
	db.Exec(`INSERT INTO questions VALUES(?, '{"1":"Correct","2":"Wrong B","3":"Wrong C","4":"Wrong D"}',1)`, question.String())
	gdb := goqu.New("sqlite3", db)
	tx, err := gdb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err = models.CopyRoomResponses(tx, participant, session); err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	if err = tx.Commit(); err != nil {
		t.Fatal(err)
	}
	var saved string
	if err = db.QueryRow(`SELECT option_order FROM user_quiz_responses`).Scan(&saved); err != nil || saved == "" {
		t.Fatalf("new host has no persisted shuffle: %q %v", saved, err)
	}
	if _, err = db.Exec(`UPDATE user_quiz_responses SET option_order='{"1":2,"2":3,"3":4,"4":1}'`); err != nil {
		t.Fatal(err)
	}

	qc := &quizSocketController{db: gdb, userQuizResponseModel: models.InitUserQuizResponseModel(gdb)}
	options := map[string]string{"1": "Correct", "2": "Wrong B", "3": "Wrong C", "4": "Wrong D"}
	for _, event := range []string{constants.EventSendQuestion, constants.EventShowScore} {
		response := &QuizSendResponse{Data: map[string]any{"id": question, "options": options, "answers": []int{1}}}
		for i := 0; i < 2; i++ { // Reconnect must reuse the exact same saved order.
			raw, err := responseForHost(qc, response, event, session.String())
			if err != nil {
				t.Fatal(err)
			}
			data := raw.(map[string]any)["data"].(map[string]any)
			var got map[string]string
			encoded, _ := json.Marshal(data["options"])
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			if got["4"] != "Correct" || got["1"] != "Wrong B" {
				t.Fatalf("wrong visible order: %v", got)
			}
			if event == constants.EventShowScore && !reflect.DeepEqual(data["answers"], []any{4}) {
				t.Fatalf("wrong revealed answer: %v", data["answers"])
			}
			if options["1"] != "Correct" {
				t.Fatal("host personalization mutated canonical broadcast")
			}
		}
	}
	translated, err := qc.userQuizResponseModel.TranslateAnswerKeys(participant, question, []int{4})
	if err != nil || !reflect.DeepEqual(translated, []int{1}) {
		t.Fatalf("displayed D must score as original A: %v, %v", translated, err)
	}
	db.Exec(`DELETE FROM user_played_quizzes`)
	response := &QuizSendResponse{Data: map[string]any{"id": question, "options": options}}
	raw, err := responseForHost(qc, response, constants.EventSendQuestion, session.String())
	if err != nil || raw != response {
		t.Fatalf("spectator host changed: %v", err)
	}
}
