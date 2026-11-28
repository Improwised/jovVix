package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// GitHub #402: the join page sends the player's name URL-encoded in the path.
// A name with a space or in a non-Latin script must join, stored as typed,
// within the 12-character username column. Set JOVVIX_TEST_DSN to run it.
func TestCreateGuestUserNamesPostgres(t *testing.T) {
	dsn := os.Getenv("JOVVIX_TEST_DSN")
	if dsn == "" {
		t.Skip("set JOVVIX_TEST_DSN to run the guest user Postgres test")
	}
	admin, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := "guest_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
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
	// The columns InsertUser writes, with username as the migrations define it.
	if _, err = db.Exec(`CREATE TABLE users(id bpchar(20) PRIMARY KEY, first_name text, last_name text,
 email text, password text, roles text, img_key text default '',
 username varchar(12) NOT NULL CONSTRAINT users_username_ukey UNIQUE)`); err != nil {
		t.Fatal(err)
	}

	ctrl, err := NewUserController(goqu.New("postgres", db), zap.NewNop(), config.AppConfig{
		Secret:    "test-secret-at-least-32-bytes-long!!",
		JWTIssuer: "test",
		Kratos:    config.KratosConfig{CookieExpirationTime: "1h"},
	})
	if err != nil {
		t.Fatal(err)
	}
	app := fiber.New()
	app.Post("/user/:"+constants.Username, ctrl.CreateGuestUser)

	join := func(name string) (int, string) {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, "/user/"+url.PathEscape(name)+"?avatar_name=Destiny", nil)
		res, err := app.Test(req, 5000)
		if err != nil {
			t.Fatal(err)
		}
		var body struct {
			Data struct {
				Username string `json:"username"`
			} `json:"data"`
		}
		_ = json.NewDecoder(res.Body).Decode(&body)
		return res.StatusCode, body.Data.Username
	}

	for _, name := range []string{
		"Phone Player", // 12 characters, 14 once encoded
		"Raj K",        // short name with a space
		"राजेश कुमार",  // Devanagari with a space
		"ઉર્વશી",       // Gujarati
		"Aditi 🎉",      // emoji
	} {
		t.Run(name, func(t *testing.T) {
			status, got := join(name)
			if status != http.StatusOK || got != name {
				t.Fatalf("first join: got %d %q, want 200 %q", status, got, name)
			}
			// A second player with the same name gets a suffixed username.
			status, got = join(name)
			if status != http.StatusOK || got == name || !utf8.ValidString(got) || utf8.RuneCountInString(got) > 12 {
				t.Fatalf("duplicate join: got %d %q, want 200 with a new valid name of at most 12 characters", status, got)
			}
		})
	}
}

// Names the username column can't hold are a 400, refused before any insert,
// so this runs without a database.
func TestCreateGuestUserRefusesBadNames(t *testing.T) {
	app := fiber.New()
	app.Post("/user/:"+constants.Username, (&UserController{logger: zap.NewNop()}).CreateGuestUser)
	for _, path := range []string{
		"/user/" + url.PathEscape("Thirteen char"),
		"/user/" + url.PathEscape("राजेश कुमार यादव"),
		"/user/bad%zzname",
	} {
		req := httptest.NewRequest(http.MethodPost, "/user/x?avatar_name=Destiny", nil)
		// app.Test sends RequestURI as written; NewRequest refuses an invalid escape.
		req.RequestURI = path + "?avatar_name=Destiny"
		res, err := app.Test(req, 5000)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("%s: got %d, want 400", path, res.StatusCode)
		}
	}
}
