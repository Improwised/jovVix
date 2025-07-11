package v1

import (
	"errors"
	"net/http"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AnalyticsBoardUserController struct {
	analyticsBoardUserModel *models.AnalyticsBoardUserModel
	presignedURLSvc         *services.PresignURLService
	logger                  *zap.Logger
}

func NewAnalyticsBoardUserController(goqu *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*AnalyticsBoardUserController, error) {
	analyticsBoardUserModel, err := models.InitAnalyticsBoardUserModel(goqu)
	if err != nil {
		return nil, err
	}

	presignedURLSvc, err := services.NewFileUploadServices(&appConfig.AWS)
	if err != nil {
		return nil, err
	}

	return &AnalyticsBoardUserController{
		analyticsBoardUserModel: &analyticsBoardUserModel,
		presignedURLSvc:         presignedURLSvc,
		logger:                  logger,
	}, nil

}

// GetAnalytics to send analytics board after quiz over to user
// swagger:route GET /v1/analytics_board/user AnalyticsBoard RequestAnalyticsBoardForUser
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
	userPlayedQuizId := ctx.Query(constants.UserPlayedQuiz)
	fc.logger.Debug("AnalyticsBoardUserController.GetAnalyticsForUser called", zap.Any("userPlayedQuizId", userPlayedQuizId))

	if userPlayedQuizId == "" || !(len(userPlayedQuizId) == 36) {
		fc.logger.Error("user played quiz is not valid")
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("user play quiz should be valid string").Error())
	}

	fc.logger.Debug("analyticsBoardUserModel.GetAnalyticsForUser called", zap.Any("userPlayedQuizId", userPlayedQuizId))
	analyticsBoardData, err := fc.analyticsBoardUserModel.GetAnalyticsForUser(userPlayedQuizId)
	if err != nil {
		fc.logger.Error("Error while getting analytics for user", zap.Error(err))
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}
	fc.logger.Debug("analyticsBoardUserModel.GetAnalyticsForUser success", zap.Any("analyticsBoardData", analyticsBoardData))

	services.ProcessAnalyticsData(analyticsBoardData, fc.presignedURLSvc, fc.logger)

	fc.logger.Debug("AnalyticsBoardUserController.GetAnalyticsForUser success", zap.Any("analyticsBoardData", analyticsBoardData))

	return utils.JSONSuccess(ctx, http.StatusOK, analyticsBoardData)
}
