package v1

import (
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func TestParseSessionIDs(t *testing.T) {
	id := uuid.NewString()
	got, ok := parseSessionIDs(" " + id + ",," + strings.ToUpper(id))
	if !ok || len(got) != 2 || got[0] != id || got[1] != id {
		t.Fatalf("got %v %v", got, ok)
	}
	if got, ok := parseSessionIDs(""); !ok || len(got) != 0 {
		t.Fatalf("empty: got %v %v", got, ok)
	}
	if _, ok := parseSessionIDs(id + ",1 OR 1=1"); ok {
		t.Fatal("accepted a value that is not a uuid")
	}
}

// TestCombinedReportPostgres runs the combined report and the reports list
// against Postgres in a throwaway schema. Set JOVVIX_TEST_DSN to run it.
func TestCombinedReportPostgres(t *testing.T) {
	dsn := os.Getenv("JOVVIX_TEST_DSN")
	if dsn == "" {
		t.Skip("set JOVVIX_TEST_DSN to run the combined report Postgres test")
	}
	admin, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := "combined_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err = admin.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec("DROP SCHEMA " + schema + " CASCADE")
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	db, err := sql.Open("postgres", dsn+sep+"search_path="+schema)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`CREATE TABLE users(id bpchar(20) PRIMARY KEY,first_name text,username text,img_key text,email text);
 CREATE TABLE quizzes(id uuid PRIMARY KEY,title text,description text,creator_id bpchar(20));
 CREATE TABLE shared_quizzes(quiz_id uuid,shared_to text);
 CREATE TABLE active_quizzes(id uuid PRIMARY KEY,admin_id bpchar(20),quiz_id uuid,activated_from timestamp,activated_to timestamp);
 CREATE TABLE active_quiz_questions(active_quiz_id uuid,question_id uuid);
 CREATE TABLE user_played_quizzes(id uuid PRIMARY KEY,user_id bpchar(20),active_quiz_id uuid,is_host bool,player_name text);
 CREATE TABLE user_quiz_responses(id uuid PRIMARY KEY,user_played_quiz_id uuid,question_id uuid,calculated_score int,response_time int);
 INSERT INTO users(id,first_name,username,img_key) VALUES('host','Host','host',''),('other','Other','other',''),('a','Asha','asha',''),('b','Bilal','bilal',''),
 ('c','Chen','chen',''),('d','Dana','dana',''),('e','Eve','eve',''),('f','Fay','fay','');`)

	quiz, otherQuiz := uuid.NewString(), uuid.NewString()
	exec(`INSERT INTO quizzes(id,title,description,creator_id) VALUES($1,'Space quiz','Batch test','host'),($2,'Other quiz','','host')`, quiz, otherQuiz)
	session := func(host, quizID, from string, ended bool) string {
		id := uuid.NewString()
		var to any
		if ended {
			to = from
		}
		exec(`INSERT INTO active_quizzes VALUES($1,$2,$3,$4,$5)`, id, host, quizID, from, to)
		exec(`INSERT INTO active_quiz_questions VALUES($1,$2),($1,$3)`, id, uuid.NewString(), uuid.NewString())
		return id
	}
	// play stores one player's two answers in a session.
	play := func(sessionID, user string, isHost bool, scores [2]int, times [2]int) {
		id := uuid.NewString()
		exec(`INSERT INTO user_played_quizzes VALUES($1,$2,$3,$4,NULL)`, id, user, sessionID, isHost)
		exec(`INSERT INTO user_quiz_responses SELECT gen_random_uuid(),$1,q.question_id,
 CASE WHEN row_number() OVER (ORDER BY q.question_id)=1 THEN $3::int ELSE $4::int END,
 CASE WHEN row_number() OVER (ORDER BY q.question_id)=1 THEN $5::int ELSE $6::int END
 FROM active_quiz_questions q WHERE q.active_quiz_id=$2`, id, sessionID, scores[0], scores[1], times[0], times[1])
	}
	s1 := session("host", quiz, "2026-11-01 10:00:00", true)
	s2 := session("host", quiz, "2026-11-02 10:00:00", true)
	running := session("host", quiz, "2026-11-03 10:00:00", false)
	foreign := session("other", quiz, "2026-11-04 10:00:00", true)
	single := session("host", otherQuiz, "2026-11-05 10:00:00", true)
	play(s1, "a", false, [2]int{50, 0}, [2]int{1000, 2000})
	play(s1, "b", false, [2]int{30, 0}, [2]int{2000, 2000})
	play(s1, "host", true, [2]int{100, 100}, [2]int{100, 100})
	play(s2, "a", false, [2]int{40, 30}, [2]int{1000, 1000})
	play(s2, "c", false, [2]int{70, 0}, [2]int{500, 500})
	play(s2, "d", false, [2]int{0, 30}, [2]int{700, 700})
	play(running, "e", false, [2]int{90, 90}, [2]int{100, 100})
	play(foreign, "f", false, [2]int{90, 90}, [2]int{100, 100})
	play(single, "a", false, [2]int{10, 10}, [2]int{100, 100})

	gdb := goqu.New("postgres", db)
	qc := &QuizController{quizModel: models.InitQuizModel(gdb), logger: zap.NewNop()}
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(constants.ContextUid, c.Get("X-User"))
		return c.Next()
	})
	app.Get("/quizzes/:"+constants.QuizId+"/combined", qc.GetCombinedReport)
	request := func(want int, user, path string) models.CombinedReport {
		t.Helper()
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Set("X-User", user)
		res, err := app.Test(req, 10000)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var body struct {
			Data json.RawMessage `json:"data"`
		}
		if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != want {
			t.Fatalf("GET %s as %s: got %d want %d: %s", path, user, res.StatusCode, want, body.Data)
		}
		var report models.CombinedReport
		if want == 200 {
			if err = json.Unmarshal(body.Data, &report); err != nil {
				t.Fatal(err)
			}
		}
		return report
	}
	type row struct {
		Name  string
		Score int
		Rank  int
		Plays int
	}
	rows := func(r models.CombinedReport) []row {
		out := []row{}
		for _, e := range r.Leaderboard {
			out = append(out, row{e.Name, e.Score, e.Rank, e.Plays})
		}
		return out
	}
	base := "/quizzes/" + quiz + "/combined"

	// All sessions: the running session, the other host's session and the host's own
	// play are left out. Asha played twice and keeps her best play.
	all := request(200, "host", base)
	if all.Title != "Space quiz" || len(all.Sessions) != 2 || all.Sessions[0].ID != s1 || all.Sessions[1].ID != s2 ||
		!all.Sessions[0].Selected || !all.Sessions[1].Selected || all.Sessions[0].Participants != 2 {
		t.Fatalf("sessions: %+v", all)
	}
	if all.Participants != 4 || all.Plays != 5 {
		t.Fatalf("participants=%d plays=%d", all.Participants, all.Plays)
	}
	want := []row{{"Chen", 70, 1, 1}, {"Asha", 70, 1, 2}, {"Dana", 30, 2, 1}, {"Bilal", 30, 2, 1}}
	if got := rows(all); !slices.Equal(got, want) {
		t.Fatalf("leaderboard %+v want %+v", got, want)
	}
	if e := all.Leaderboard[1]; e.Correct != 2 || e.Questions != 2 || e.SessionID != s2 || e.ResponseTime != 2000 {
		t.Fatalf("Asha's entry: %+v", e)
	}

	// One session selected: only its players, the other session is listed but not selected.
	one := request(200, "host", base+"?sessions="+s1)
	if one.Sessions[1].Selected || !one.Sessions[0].Selected {
		t.Fatalf("selection: %+v", one.Sessions)
	}
	if got := rows(one); !slices.Equal(got, []row{{"Asha", 50, 1, 1}, {"Bilal", 30, 2, 1}}) {
		t.Fatalf("one session: %+v", got)
	}

	// Another host's session, a running session and a session of another quiz are refused.
	request(400, "host", base+"?sessions="+s1+","+foreign)
	request(400, "host", base+"?sessions="+running)
	request(400, "host", base+"?sessions="+single)
	request(400, "host", base+"?sessions=nope")
	request(400, "host", "/quizzes/nope/combined")
	// The other host sees only their own session of the same quiz.
	theirs := request(200, "other", base)
	if len(theirs.Sessions) != 1 || theirs.Sessions[0].ID != foreign || theirs.Participants != 1 {
		t.Fatalf("other host: %+v", theirs)
	}
	// A user who never hosted the quiz gets 404, not the host's data.
	request(404, "a", base)

	// The reports list names the quiz and how many ended sessions it has.
	list, count, err := qc.quizModel.ListQuizzesAnalysis("", "", "", "", "host", 1)
	if err != nil || count != 3 {
		t.Fatalf("list: count=%d err=%v", count, err)
	}
	sessions := map[string]int{}
	for _, q := range list {
		sessions[q.ID.String()] = q.Sessions
		if q.ID.String() == single && q.QuizID.String() != otherQuiz {
			t.Fatalf("quiz id: %+v", q)
		}
	}
	if sessions[s1] != 2 || sessions[s2] != 2 || sessions[single] != 1 {
		t.Fatalf("session counts: %v", sessions)
	}
}
