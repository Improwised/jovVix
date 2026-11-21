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

const reportSession = "6f1c2a3e-0000-4000-8000-000000000001"

// GitHub #38: the creator of a shared quiz reads reports of sessions someone else hosted.
func TestVerifyReportAccess(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`CREATE TABLE quizzes(id text PRIMARY KEY, creator_id text);
CREATE TABLE active_quizzes(id text PRIMARY KEY, quiz_id text, admin_id text);
INSERT INTO quizzes VALUES('quiz-1', 'owner');
INSERT INTO active_quizzes VALUES('` + reportSession + `', 'quiz-1', 'host');`); err != nil {
		t.Fatal(err)
	}
	m := Middleware{Db: goqu.New("sqlite3", db), Logger: zap.NewNop()}

	for user, want := range map[string]int{
		"host":     http.StatusOK,
		"owner":    http.StatusOK,
		"stranger": http.StatusForbidden,
	} {
		app := fiber.New()
		app.Get("/report", func(c *fiber.Ctx) error {
			c.Locals(constants.ContextUid, user)
			return c.Next()
		}, m.VerifyReportAccess, func(c *fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

		res, err := app.Test(httptest.NewRequest(http.MethodGet, "/report?"+constants.ActiveQuizId+"="+reportSession, nil))
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != want {
			t.Errorf("%s: status %d, want %d", user, res.StatusCode, want)
		}
	}
}
