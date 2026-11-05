package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func TestHostParticipationPostgres(t *testing.T) {
	dsn := os.Getenv("HOST_PARTICIPATION_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set HOST_PARTICIPATION_TEST_DATABASE_URL for database test")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	schema := "host_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err = db.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	defer db.Exec("DROP SCHEMA " + schema + " CASCADE")
	if _, err = db.Exec("SET search_path TO " + schema); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE active_quizzes (id uuid PRIMARY KEY, admin_id text NOT NULL, is_active boolean NOT NULL, current_question uuid, activated_to timestamptz, is_question_active boolean);
 CREATE TABLE active_quiz_questions (active_quiz_id uuid,question_id uuid);
 CREATE TABLE users (id text, first_name text, last_name text, username text, img_key text);
 INSERT INTO users VALUES ('host','Host','One','host','avatar'), ('player','Host','Two','player','avatar');
 CREATE TABLE user_played_quizzes (id uuid PRIMARY KEY,user_id text,active_quiz_id uuid,is_host boolean, player_name text);
 CREATE TABLE user_quiz_responses (id uuid PRIMARY KEY,question_id uuid,user_played_quiz_id uuid,option_order jsonb,answers jsonb,is_attend boolean DEFAULT false);
 INSERT INTO active_quizzes VALUES ('00000000-0000-0000-0000-000000000001','host',true,NULL,NULL,NULL);
 INSERT INTO active_quiz_questions VALUES ('00000000-0000-0000-0000-000000000001','00000000-0000-0000-0000-000000000011');`)
	if err != nil {
		t.Fatal(err)
	}
	qc := &quizSocketController{db: goqu.New("postgres", db)}
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(constants.ContextUser, models.User{ID: c.Get("X-User"), FirstName: "Host", Username: "host"})
		return c.Next()
	})
	app.Get("/participation", qc.HostParticipation)
	app.Put("/participation", qc.HostParticipation)
	request := func(method, user, body string) (int, map[string]any) {
		t.Helper()
		req := httptest.NewRequest(method, "/participation?active_quiz_id=00000000-0000-0000-0000-000000000001", strings.NewReader(body))
		req.Header.Set("X-User", user)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var payload map[string]any
		if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		return res.StatusCode, payload
	}
	session := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	if ok, err := qc.hasHostedParticipants(session); err != nil || ok {
		t.Fatalf("empty start: %v %v", ok, err)
	}
	if code, _ := request("PUT", "other", `{"enabled":true}`); code != 403 {
		t.Fatalf("other host status: %d", code)
	}
	playerModel := models.InitUserPlayedQuizModel(goqu.New("postgres", db))
	_, _, _, _, err = playerModel.JoinQuiz("player", session, "Host")
	if err != nil {
		t.Fatal(err)
	}
	code, conflict := request("PUT", "host", `{"enabled":true}`)
	if code != 409 || len(conflict["data"].(map[string]any)["suggestions"].([]any)) == 0 {
		t.Fatalf("expected name conflict and suggestions: %d %v", code, conflict)
	}
	code, state := request("PUT", "host", `{"enabled":true,"player_name":"HostTwo"}`)
	if code != 200 || state["data"].(map[string]any)["name"] != "HostTwo" {
		t.Fatalf("chosen name: %d %v", code, state)
	}
	code, state = request("GET", "host", "")
	if code != 200 || state["data"].(map[string]any)["name"] != "HostTwo" {
		t.Fatalf("saved name: %d %v", code, state)
	}
	if code, _ = request("PUT", "host", `{"enabled":false}`); code != 200 {
		t.Fatalf("disable chosen name: %d", code)
	}
	if _, err = db.Exec(`DELETE FROM user_played_quizzes WHERE user_id='player'`); err != nil {
		t.Fatal(err)
	}
	code, state = request("PUT", "host", `{"enabled":true}`)
	if _, _, _, _, err = playerModel.JoinQuiz("player", session, "hOsT"); !errors.Is(err, models.ErrPlayerNameTaken) {
		t.Fatalf("host-first duplicate was accepted: %v", err)
	}
	if code != 200 || state["data"].(map[string]any)["enabled"] != true {
		t.Fatalf("enable: %d %v", code, state)
	}
	id := state["data"].(map[string]any)["user_played_quiz"]
	code, state = request("PUT", "host", `{"enabled":true}`)
	if code != 200 || state["data"].(map[string]any)["user_played_quiz"] != id {
		t.Fatalf("idempotent enable: %v", state)
	}
	var count int
	if err = db.QueryRow("SELECT count(*) FROM user_quiz_responses").Scan(&count); err != nil || count != 1 {
		t.Fatalf("responses %d: %v", count, err)
	}
	if code, _ = request("PUT", "host", `{"enabled":false}`); code != 200 {
		t.Fatalf("disable: %d", code)
	}
	if err = db.QueryRow("SELECT count(*) FROM user_played_quizzes").Scan(&count); err != nil || count != 0 {
		t.Fatalf("disable rows %d: %v", count, err)
	}
	if code, _ = request("PUT", "host", `{"enabled":true}`); code != 200 {
		t.Fatalf("enable again: %d", code)
	}
	if ok, err := qc.hasHostedParticipants(session); err != nil || !ok {
		t.Fatalf("solo start: %v %v", ok, err)
	}
	// Starting locks participation immediately, even before the first question.
	code, state = request("GET", "host", "")
	if code != 200 || state["data"].(map[string]any)["started"] != true {
		t.Fatalf("countdown state: %d %v", code, state)
	}
	if code, _ = request("PUT", "host", `{"enabled":false}`); code != 409 {
		t.Fatalf("participation changed during countdown: %d", code)
	}
	var questionActive sql.NullBool
	if err = db.QueryRow("SELECT is_question_active FROM active_quizzes WHERE id=$1", session).Scan(&questionActive); err != nil || !questionActive.Valid || questionActive.Bool {
		t.Fatalf("countdown marker: %v %v", questionActive, err)
	}

	_, err = db.Exec(`UPDATE active_quizzes SET current_question='00000000-0000-0000-0000-000000000011', is_question_active=true; UPDATE user_quiz_responses SET answers='[1]', is_attend=true`)
	if err != nil {
		t.Fatal(err)
	}
	code, state = request("GET", "host", "")
	restored := state["data"].(map[string]any)
	if code != 200 || len(restored["answer"].([]any)) != 1 || len(restored["submitted"].([]any)) != 1 || restored["started"] != true {
		t.Fatalf("restore: %d %v", code, state)
	}
	if code, _ = request("PUT", "host", `{"enabled":false}`); code != 409 {
		t.Fatalf("after start: %d", code)
	}
}
