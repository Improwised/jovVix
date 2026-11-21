package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GitHub #77: an expired login gets 401 and a message a person understands, and a Kratos outage
// never answers with status 0.
func TestKratosAuthenticatedSessionFailures(t *testing.T) {
	kratosStatus := http.StatusUnauthorized
	kratos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(kratosStatus)
	}))
	defer kratos.Close()
	down := httptest.NewServer(http.NotFoundHandler())
	downURL := down.URL
	down.Close()

	cases := []struct {
		name, baseURL, cookie string
		kratosStatus          int
		wantStatus            int
		wantMessage           string
	}{
		{"no session cookie", kratos.URL, "", 0, http.StatusUnauthorized, constants.ErrSessionExpired},
		{"expired session", kratos.URL, "old", http.StatusUnauthorized, http.StatusUnauthorized, constants.ErrSessionExpired},
		{"kratos error", kratos.URL, "s", http.StatusInternalServerError, http.StatusInternalServerError, constants.ErrKratosAuth},
		{"kratos unreachable", downURL, "s", 0, http.StatusBadGateway, constants.ErrKratosAuth},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kratosStatus = tc.kratosStatus
			m := Middleware{Config: config.AppConfig{Kratos: config.KratosConfig{BaseUrl: tc.baseURL}}, Logger: zap.NewNop()}
			app := fiber.New()
			app.Get("/admin", m.KratosAuthenticated, func(c *fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: "ory_kratos_session", Value: tc.cookie})
			}
			res, err := app.Test(req, 5000)
			if err != nil {
				t.Fatal(err)
			}
			var body struct {
				Message string `json:"message"`
			}
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tc.wantStatus || body.Message != tc.wantMessage {
				t.Fatalf("got %d %q, want %d %q", res.StatusCode, body.Message, tc.wantStatus, tc.wantMessage)
			}
		})
	}
}
