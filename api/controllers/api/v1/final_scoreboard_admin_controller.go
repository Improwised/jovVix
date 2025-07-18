package v1

import (
	"errors"
	"net/http"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type FinalScoreBoardAdminController struct {
	finalScoreBoardAdminModel *models.FinalScoreBoardAdminModel
	logger                    *zap.Logger
}

func NewFinalScoreBoardAdminController(goqu *goqu.Database, logger *zap.Logger) (*FinalScoreBoardAdminController, error) {
	finalScoreBoardAdminModel, err := models.InitFinalScoreBoardAdminModel(goqu)
	if err != nil {
		return nil, err
	}

	return &FinalScoreBoardAdminController{
		finalScoreBoardAdminModel: &finalScoreBoardAdminModel,
		logger:                    logger,
	}, nil

}

// GetScore to send final score after quiz over to admin
// swagger:route GET /v1/final_score/admin FinalScore RequestFinalScoreForAdmin
//
// Get a finalScore details for admin.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseFinalScoreForAdmin
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (fc *FinalScoreBoardAdminController) GetScoreForAdmin(ctx *fiber.Ctx) error {

	var activeQuizId = ctx.Query(constants.ActiveQuizId)

	if !(activeQuizId != "" && len(activeQuizId) == 36) {
		fc.logger.Debug("active quiz id is not valid - either empty string or it is not 36 characters long")
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("user play quiz should be valid string").Error())
	}

	finalScoreBoardData, err := fc.finalScoreBoardAdminModel.GetScoreForAdmin(activeQuizId)
	if err != nil {
		fc.logger.Error("Error while getting final scoreboard for admin", zap.Error(err))
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}

	return utils.JSONSuccess(ctx, http.StatusOK, finalScoreBoardData)
}
