package v1

import (
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func TestQuizRoomsPostgres(t *testing.T) {
	dsn := os.Getenv("ROOM_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set ROOM_TEST_DATABASE_URL for room transaction and migration tests")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	schema := "room_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err = db.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	defer db.Exec("DROP SCHEMA " + schema + " CASCADE")
	// Every connection uses the isolated schema, including concurrent requests.
	isolated, err := sql.Open("postgres", dsn+"&search_path="+schema)
	if err != nil {
		t.Fatal(err)
	}
	defer isolated.Close()
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := isolated.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`CREATE TABLE users(id bpchar(20) PRIMARY KEY,first_name text,last_name text);
 CREATE TABLE quizzes(id uuid PRIMARY KEY,title text,is_public bool,creator_id bpchar(20));
 CREATE TABLE shared_quizzes(quiz_id uuid,shared_to text);
 CREATE TABLE questions(id uuid PRIMARY KEY,options json DEFAULT '{"1":"Yes"}',type int DEFAULT 1);
 CREATE TABLE user_quiz_responses(id uuid PRIMARY KEY,question_id uuid,user_played_quiz_id uuid,option_order json,answers json,is_attend bool DEFAULT false);
 CREATE TABLE quiz_questions(id uuid,quiz_id uuid,question_id uuid,next_question uuid,created_at timestamp DEFAULT now());
 CREATE TABLE active_quizzes(id uuid PRIMARY KEY,admin_id bpchar(20),quiz_id uuid,title text,invitation_code integer,is_active bool DEFAULT false,activated_from timestamp,activated_to timestamp,current_question uuid,is_question_active bool,is_solo bool DEFAULT false);
 CREATE TABLE user_played_quizzes(id uuid PRIMARY KEY,user_id bpchar(20),active_quiz_id uuid,is_host bool,player_name text);
 CREATE TABLE active_quiz_questions(id uuid,question_id uuid,active_quiz_id uuid,next_question uuid,order_no int);
 INSERT INTO users VALUES('host','Host',''),('one','Alpha',''),('two','Beta',''),('new','Gamma','');`)
	exec(`ALTER TABLE users ADD COLUMN username text DEFAULT ''; ALTER TABLE users ADD COLUMN img_key text; ALTER TABLE user_quiz_responses ADD COLUMN calculated_score int DEFAULT 0;`)
	session, quiz, empty, private, question := uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()
	exec(`INSERT INTO quizzes VALUES($1,'Public',true,'host'),($2,'Empty',true,'host'),($3,'Private',false,'two')`, quiz, empty, private)
	exec(`INSERT INTO questions(id) VALUES($1)`, question)
	exec(`INSERT INTO quiz_questions(id,quiz_id,question_id) VALUES($1,$2,$3)`, uuid.New(), quiz, question)
	exec(`INSERT INTO active_quizzes(id,admin_id,quiz_id,title,invitation_code,is_active) VALUES($1,'host',$2,'Old round',123456,true)`, session, quiz)
	exec(`INSERT INTO user_played_quizzes VALUES($1,'one',$2,false,'Alpha')`, uuid.New(), session)
	migration, err := os.ReadFile("../../../database/migrations/20261105120000_quiz_rooms.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	exec(string(migration))
	var room uuid.UUID
	if err = isolated.QueryRow(`SELECT room_id FROM active_quizzes WHERE id=$1`, session).Scan(&room); err != nil {
		t.Fatal(err)
	}
	qc := &quizSocketController{db: goqu.New("postgres", isolated)}
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(constants.ContextUser, models.User{ID: c.Get("X-User"), FirstName: "Player", Username: c.Get("X-User")})
		return c.Next()
	})
	app.Put("/host", qc.HostParticipation)
	app.Get("/room-reports", qc.RoomReports)
	app.Get("/rooms/:room", qc.RoomState)
	app.Get("/rooms/:room/report", qc.RoomReport)
	app.Post("/join/:code", qc.JoinRoom)
	app.Post("/rooms/:room/leave", qc.LeaveRoom)
	app.Post("/rooms/:room/next", qc.NextRoomQuiz)
	app.Post("/rooms/:room/close", qc.CloseRoom)
	app.Get("/rooms/:room/quizzes", qc.RoomQuizzes)
	request := func(method, path, user, body string) (int, map[string]any) {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("X-User", user)
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, 10000)
		if err != nil {
			t.Error(err)
			return 0, nil
		}
		defer res.Body.Close()
		var payload map[string]any
		if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		return res.StatusCode, payload
	}
	path := "/rooms/" + room.String()
	check := func(want int, method, path, user, body string) map[string]any {
		t.Helper()
		got, payload := request(method, path, user, body)
		if got != want {
			t.Fatalf("%s %s: got %d want %d: %v", method, path, got, want, payload)
		}
		return payload
	}
	reportList := check(200, "GET", "/room-reports", "host", "")["data"].([]any)
	row := reportList[0].(map[string]any)
	if _, err := time.Parse(time.RFC3339Nano, row["started"].(string)); err != nil {
		t.Fatalf("room start has no timezone: %v", err)
	}
	if row["finished"] != nil {
		t.Fatalf("open room has end date: %v", row["finished"])
	}
	// Migration preserves the live code and the existing member.
	check(200, "GET", path, "one", "")
	check(403, "GET", path, "two", "")
	check(403, "GET", path+"/quizzes", "one", "")
	check(403, "POST", path+"/close", "one", "")
	check(409, "POST", path+"/close", "host", "")
	check(200, "POST", "/join/123456", "two", `{"player_name":"Beta"}`)
	check(409, "POST", "/join/123456", "new", `{"player_name":"Alpha"}`)
	check(200, "GET", path+"/report", "two", "")
	check(200, "POST", path+"/leave", "two", "")
	check(403, "GET", path+"/report", "two", "")
	check(200, "GET", path+"/report", "host", "")
	check(403, "GET", path+"/report", "new", "")
	check(403, "GET", path, "two", "")
	// Names stay reserved for a returning member across rounds.
	check(409, "POST", "/join/123456", "new", `{"player_name":"Beta"}`)
	check(200, "POST", "/join/123456", "two", `{"player_name":"Changed"}`)
	var name string
	if err = isolated.QueryRow(`SELECT player_name FROM quiz_room_members WHERE room_id=$1 AND user_id='two'`, room).Scan(&name); err != nil || name != "Beta" {
		t.Fatalf("returning name %q: %v", name, err)
	}
	exec(`UPDATE active_quizzes SET is_active=false,invitation_code=NULL,activated_to=now() WHERE id=$1`, session)
	check(200, "POST", "/join/123456", "new", `{"player_name":"Gamma"}`)
	body := func(q uuid.UUID) string {
		return `{"quiz_id":"` + q.String() + `","expected_session_id":"` + session.String() + `"}`
	}
	check(403, "POST", path+"/next", "one", body(quiz))
	check(403, "POST", path+"/next", "host", body(private))
	check(400, "POST", path+"/next", "host", body(empty))
	var count int
	isolated.QueryRow(`SELECT count(*) FROM active_quizzes`).Scan(&count)
	if count != 1 {
		t.Fatalf("failed selection leaked rounds: %d", count)
	}
	codes := make(chan int, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); code, _ := request("POST", path+"/next", "host", body(quiz)); codes <- code }()
	}
	wg.Wait()
	close(codes)
	counts := map[int]int{}
	for code := range codes {
		counts[code]++
	}
	if counts[201] != 1 || counts[409] != 1 {
		t.Fatalf("concurrent transitions: %v", counts)
	}
	var current uuid.UUID
	isolated.QueryRow(`SELECT current_session_id FROM quiz_rooms WHERE id=$1`, room).Scan(&current)
	if current == session {
		t.Fatal("round was reused")
	}
	check(200, "GET", path, "one", "")
	check(409, "POST", path+"/close", "host", "")
	// An offline retained member has no participant row in this new round,
	// but their room name must still be protected from host participation.
	hostPath := "/host?active_quiz_id=" + current.String()
	check(409, "PUT", hostPath, "host", `{"enabled":true,"player_name":"Alpha"}`)
	check(200, "PUT", hostPath, "host", `{"enabled":true,"player_name":"HostPlayer"}`)
	exec(`UPDATE active_quizzes SET is_active=false WHERE id=$1`, current)
	check(201, "POST", path+"/next", "host", `{"quiz_id":"`+quiz.String()+`","expected_session_id":"`+current.String()+`"}`)
	isolated.QueryRow(`SELECT current_session_id FROM quiz_rooms WHERE id=$1`, room).Scan(&current)
	var hostCount int
	isolated.QueryRow(`SELECT count(*) FROM user_played_quizzes WHERE active_quiz_id=$1 AND user_id='host' AND player_name='HostPlayer'`, current).Scan(&hostCount)
	if hostCount != 1 {
		t.Fatal("playing host was not carried to the next round")
	}

	exec(`UPDATE active_quizzes SET is_active=false WHERE id=$1`, current)
	check(200, "POST", path+"/close", "host", "")
	check(200, "POST", path+"/leave", "two", "")
	check(200, "GET", path+"/report", "two", "")
	reportList = check(200, "GET", "/room-reports", "host", "")["data"].([]any)
	row = reportList[0].(map[string]any)
	if _, err := time.Parse(time.RFC3339Nano, row["finished"].(string)); err != nil {
		t.Fatalf("room end has no timezone: %v", err)
	}
	check(200, "POST", path+"/close", "host", "")
	check(410, "POST", "/join/123456", "new", `{"player_name":"Gamma"}`)
	check(410, "POST", path+"/next", "host", body(quiz))
	exec(`UPDATE quiz_rooms SET closed_at=NULL,expires_at=now()-interval '1 second' WHERE id=$1`, room)
	check(410, "POST", "/join/123456", "new", `{"player_name":"Gamma"}`)
	// Reversible schema change leaves historical quiz sessions intact.
	down, err := os.ReadFile("../../../database/migrations/20261105120000_quiz_rooms.down.sql")
	if err != nil {
		t.Fatal(err)
	}
	exec(string(down))
	isolated.QueryRow(`SELECT count(*) FROM active_quizzes`).Scan(&count)
	if count != 3 {
		t.Fatalf("historical sessions lost: %d", count)
	}
}
