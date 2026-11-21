package middlewares

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	fiber "github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const (
	sharedSession = "6f1c2a3e-0000-4000-8000-000000000001" // owner's quiz, hosted by a colleague it was shared with
	publicSession = "6f1c2a3e-0000-4000-8000-000000000002" // owner's public quiz, hosted by a stranger
)

// GitHub #38: the creator of a shared quiz reads reports of sessions that the people it was
// shared with hosted, but not of sessions that anyone else hosted from a public quiz.
func TestVerifyReportAccess(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`CREATE TABLE users(id text PRIMARY KEY, email text);
CREATE TABLE quizzes(id text PRIMARY KEY, creator_id text, is_public boolean);
CREATE TABLE active_quizzes(id text PRIMARY KEY, quiz_id text, admin_id text);
CREATE TABLE shared_quizzes(quiz_id text, shared_to text);
INSERT INTO users VALUES('owner', 'owner@example.com'), ('colleague', 'colleague@example.com'), ('stranger', 'stranger@example.com');
INSERT INTO quizzes VALUES('quiz-1', 'owner', 1), ('quiz-2', 'owner', 0);
-- The stranger holds a share of another quiz only: that must not open quiz-1's sessions.
INSERT INTO shared_quizzes VALUES('quiz-1', 'colleague@example.com'), ('quiz-2', 'stranger@example.com');
INSERT INTO active_quizzes VALUES('` + sharedSession + `', 'quiz-1', 'colleague'), ('` + publicSession + `', 'quiz-1', 'stranger');`); err != nil {
		t.Fatal(err)
	}
	m := Middleware{Db: goqu.New("sqlite3", db), Logger: zap.NewNop()}

	cases := []struct {
		user, session string
		want          int
	}{
		{"colleague", sharedSession, http.StatusOK}, // the host
		{"owner", sharedSession, http.StatusOK},     // the creator, host has a share
		{"stranger", sharedSession, http.StatusForbidden},
		{"stranger", publicSession, http.StatusOK},     // the host
		{"owner", publicSession, http.StatusForbidden}, // the creator, host has no share
		{"colleague", publicSession, http.StatusForbidden},
	}
	for _, tc := range cases {
		app := fiber.New()
		app.Get("/report", func(c *fiber.Ctx) error {
			c.Locals(constants.ContextUid, tc.user)
			return c.Next()
		}, m.VerifyReportAccess, func(c *fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

		res, err := app.Test(httptest.NewRequest(http.MethodGet, "/report?"+constants.ActiveQuizId+"="+tc.session, nil))
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != tc.want {
			t.Errorf("%s on %s: status %d, want %d", tc.user, tc.session, res.StatusCode, tc.want)
		}
	}
}
