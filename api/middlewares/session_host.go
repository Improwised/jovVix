package middlewares

import (
	"net/http"

	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// VerifySessionHost allows the request only when the authenticated user hosted
// the quiz session named by the active_quiz_id route param or query string.
// It guards host-only views (reports, admin scoreboard and analytics), which
// would otherwise expose every player's answers for any session id.
func (m *Middleware) VerifySessionHost(c *fiber.Ctx) error {
	activeQuizId := c.Params(constants.ActiveQuizId)
	if activeQuizId == "" {
		activeQuizId = c.Query(constants.ActiveQuizId)
	}

	if _, err := uuid.Parse(activeQuizId); err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "active quiz id should be valid string")
	}

	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	if userId == "" {
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrUnauthenticated)
	}

	isHost, err := m.Db.From(models.ActiveQuizzesTable).
		Where(goqu.Ex{"id": activeQuizId, "admin_id": userId}).
		Select(goqu.L("1")).
		Limit(1).
		ScanVal(new(int))
	if err != nil {
		m.Logger.Error("error while verifying session host", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthorized)
	}

	if !isHost {
		return utils.JSONError(c, http.StatusForbidden, constants.ErrUnauthorized)
	}

	return c.Next()
}
