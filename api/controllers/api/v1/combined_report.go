package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// parseSessionIDs splits the comma-separated sessions query value. It rejects
// anything that is not a UUID, so the database only sees valid ids.
func parseSessionIDs(raw string) ([]string, bool) {
	ids := []string{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := uuid.Parse(part)
		if err != nil {
			return nil, false
		}
		ids = append(ids, id.String())
	}
	return ids, true
}

// GetCombinedReport merges the leaderboards of several sessions of one quiz.
// swagger:route GET /v1/admin/reports/quizzes/{quiz_id}/combined Reports RequestGetCombinedReport
//
// Get one leaderboard across the selected ended sessions of a quiz that the user hosted.
// Without the sessions query, every ended session is selected.
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseGetCombinedReport
//	     400: GenericResFailNotFound
//	     404: GenericResFailNotFound
//		  500: GenericResError
func (qc *QuizController) GetCombinedReport(c *fiber.Ctx) error {
	quizID, err := uuid.Parse(c.Params(constants.QuizId))
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, "quiz id should be a valid uuid")
	}
	selected, ok := parseSessionIDs(c.Query("sessions"))
	if !ok {
		return utils.JSONFail(c, http.StatusBadRequest, "sessions should be comma-separated session ids")
	}
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	report, err := qc.quizModel.GetCombinedReport(quizID.String(), userID, selected)
	var unknown models.ErrUnknownSession
	switch {
	case errors.As(err, &models.ErrNoSessions{}):
		return utils.JSONFail(c, http.StatusNotFound, err.Error())
	case errors.As(err, &unknown):
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	case err != nil:
		qc.logger.Error("error while building combined report", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, "unable to build the combined report")
	}
	return utils.JSONSuccess(c, http.StatusOK, report)
}
