package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// GitHub #77: /kratos/whoami says the session expired in words, and a Kratos outage is a 502.
func TestGetRegisteredUserSessionFailures(t *testing.T) {
	kratos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer kratos.Close()
	down := httptest.NewServer(http.NotFoundHandler())
	downURL := down.URL
	down.Close()

	cases := []struct {
		name, baseURL, cookie string
		wantStatus            int
		wantMessage           string
	}{
		{"no session cookie", kratos.URL, "", http.StatusUnauthorized, constants.ErrSessionExpired},
		{"expired session", kratos.URL, "old", http.StatusUnauthorized, constants.ErrSessionExpired},
		{"kratos unreachable", downURL, "s", http.StatusBadGateway, constants.ErrKratosAuth},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := &AuthController{logger: zap.NewNop(), config: config.AppConfig{Kratos: config.KratosConfig{BaseUrl: tc.baseURL}}}
			app := fiber.New()
			app.Get("/whoami", ctrl.GetRegisteredUser)
			req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
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
