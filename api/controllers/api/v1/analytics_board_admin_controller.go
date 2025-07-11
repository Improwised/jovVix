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

type AnalyticsBoardAdminController struct {
	AnalyticsBoardAdminModel *models.AnalyticsBoardAdminModel
	presignedURLSvc          *services.PresignURLService
	logger                   *zap.Logger
}

func NewAnalyticsBoardAdminController(goqu *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*AnalyticsBoardAdminController, error) {
	analyticsBoardAdminModel, err := models.InitAnalyticsBoardAdminModel(goqu)
	if err != nil {
		return nil, err
	}

	presignedURLSvc, err := services.NewFileUploadServices(&appConfig.AWS)
	if err != nil {
		return nil, err
	}

	return &AnalyticsBoardAdminController{
		AnalyticsBoardAdminModel: &analyticsBoardAdminModel,
		presignedURLSvc:          presignedURLSvc,
		logger:                   logger,
	}, nil

}

// GetAnalytics to send final score after quiz over to admin
// swagger:route GET /v1/analytics_board/admin AnalyticsBoard RequestAnalyticsBoardForAdmin
//
// Get a analyticsboard details for admin.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseAnalyticsBoardForAdmin
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (fc *AnalyticsBoardAdminController) GetAnalyticsForAdmin(ctx *fiber.Ctx) error {

	var activeQuizId = ctx.Query(constants.ActiveQuizId, "")

	if activeQuizId == "" || !(len(activeQuizId) == 36) {
		fc.logger.Error("active quiz id is not valid")
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("active quiz id should be valid string").Error())
	}

	analyticsBoardData, err := fc.AnalyticsBoardAdminModel.GetAnalyticsForAdmin(activeQuizId)
	if err != nil {
		fc.logger.Error("Error while getting analytics for admin", zap.Error(err))
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}

	services.ProcessAnalyticsData(analyticsBoardData, fc.presignedURLSvc, fc.logger)

	return utils.JSONSuccess(ctx, http.StatusOK, analyticsBoardData)

}
