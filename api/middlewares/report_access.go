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

// VerifyReportAccess allows the request only when the authenticated user hosted
// the quiz session named by the active_quiz_id route param or query string, or
// created its quiz and shared it with the host (see models.ReportVisibleTo).
// It guards read-only report views (reports, admin scoreboard and analytics),
// which would otherwise expose every player's answers for any session id.
func (m *Middleware) VerifyReportAccess(c *fiber.Ctx) error {
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

	allowed, err := m.Db.From(goqu.T(models.ActiveQuizzesTable).As("aq")).
		LeftJoin(goqu.T(models.QuizzesTable).As("q"), goqu.On(goqu.Ex{"q.id": goqu.I("aq.quiz_id")})).
		Where(
			goqu.I("aq.id").Eq(activeQuizId),
			models.ReportVisibleTo(userId),
		).
		Select(goqu.L("1")).
		Limit(1).
		ScanVal(new(int))
	if err != nil {
		m.Logger.Error("error while verifying report access", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthorized)
	}

	if !allowed {
		return utils.JSONError(c, http.StatusForbidden, constants.ErrUnauthorized)
	}

	return c.Next()
}
