package v1

import (
	"errors"
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AnalyticsBoardUserController struct {
	analyticsBoardUserModel  *models.AnalyticsBoardUserModel
	logger                    *zap.Logger
	event                     *events.Events
}

func NewAnalyticsBoardUserController(goqu *goqu.Database, logger *zap.Logger, event *events.Events) (*AnalyticsBoardUserController, error) {
	analyticsBoardUserModel, err := models.InitAnalyticsBoardUserModel(goqu)
	if err != nil {
		return nil, err
	}

	return &AnalyticsBoardUserController{
		analyticsBoardUserModel: &analyticsBoardUserModel,
		logger:                    logger,
		event:                     event,
	}, nil

}

// GetAnalytics to send analytics board after quiz over to user
// swagger:route GET /v1/analytics_board/user AnalyticsBoardForUser RequestAnalyticsBoardForUser
//
// Get a analyticsboard details for user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseAnalyticsBoardForUser
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (fc *AnalyticsBoardUserController) GetAnalyticsForUser(ctx *fiber.Ctx) error {

	var userPlayedQuizId = ctx.Cookies(constants.UserPlayedQuiz)

	if userPlayedQuizId == "" || !(len(userPlayedQuizId) == 36) {
		fc.logger.Error("user played quiz is not valid")
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("user play quiz should be valid string").Error())
	}
	analyticsBoardData, err := fc.analyticsBoardUserModel.GetAnalyticsForUser(userPlayedQuizId)
	if err != nil {
		fc.logger.Error("Error while getting analytics for user", zap.Error(err))
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}

	return utils.JSONSuccess(ctx, http.StatusOK, analyticsBoardData)

}
