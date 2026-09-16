package middlewares

import (
	"testing"

	"github.com/Improwised/jovvix/api/config"
)

func TestAIRateLimitTier(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"/api/v1/ai/questions/generate", "gen"},
		{"/api/v1/ai/quizzes", "gen"},
		{"/api/v1/ai/test", "gen"},
		{"/api/v1/ai/models", "gen"},
		{"/api/v1/ai/status", "meta"},
		{"/api/v1/ai/settings", "meta"},
		{"/api/v1/ai/settings/unlock", "meta"},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			if got := aiRateLimitTier(tc.path); got != tc.want {
				t.Fatalf("aiRateLimitTier(%q) = %q, want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestAIRateLimitBounds(t *testing.T) {
	m := &Middleware{Config: config.AppConfig{AI: config.AIConfig{GenRateLimit: 7, GenRateWindow: 30}}}

	limit, window := m.aiRateLimitBounds("gen")
	if limit != 7 {
		t.Fatalf("gen limit = %d, want 7", limit)
	}
	if window.Seconds() != 30 {
		t.Fatalf("gen window = %v, want 30s", window)
	}
}
