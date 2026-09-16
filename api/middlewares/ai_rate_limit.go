package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// aiRateLimitScript increments the counter and, on the first increment, sets
// its expiry. Both run atomically so a counter can never be created without a
// TTL (which would otherwise rate-limit a user indefinitely).
var aiRateLimitScript = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
	redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return count
`)

// AIRateLimit bounds the server-side outbound calls the AI feature makes on
// behalf of a caller. Generation/test/models routes get the stricter "gen"
// tier because they trigger provider requests; settings/status/unlock are
// cheap local reads and use the "meta" tier.
func (m *Middleware) AIRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if m.rateLimiter == nil {
			return c.Next()
		}

		tier := aiRateLimitTier(c.Path())
		limit, window := m.aiRateLimitBounds(tier)

		userID, ok := c.Locals(constants.ContextUid).(string)
		if !ok || userID == "" {
			userID = c.IP()
		}

		key := "ai:rl:" + tier + ":" + userID

		count, err := aiRateLimitScript.Run(
			m.rateLimiter.Ctx,
			m.rateLimiter.Client,
			[]string{key},
			int64(window.Seconds()),
		).Int64()
		if err != nil {
			m.Logger.Error("ai rate limit incr failed", zap.Error(err))
			return c.Next()
		}

		if count > int64(limit) {
			return utils.JSONError(c, http.StatusTooManyRequests, constants.ErrAIRateLimited)
		}

		return c.Next()
	}
}

func aiRateLimitTier(path string) string {
	switch {
	case strings.Contains(path, "/generate"),
		strings.Contains(path, "/test"),
		strings.Contains(path, "/models"),
		strings.HasSuffix(path, "/quizzes"):
		return "gen"
	default:
		return "meta"
	}
}

func (m *Middleware) aiRateLimitBounds(tier string) (int, time.Duration) {
	if tier == "gen" {
		return m.Config.AI.GenerationRateLimit(), time.Duration(m.Config.AI.GenerationRateWindow()) * time.Second
	}
	return m.Config.AI.MetaRateLimitValue(), time.Duration(m.Config.AI.MetaRateWindowValue()) * time.Second
}
