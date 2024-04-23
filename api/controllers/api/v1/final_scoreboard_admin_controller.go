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

type FinalScoreBoardAdminController struct {
	finalScoreBoardAdminModel *models.FinalScoreBoardAdminModel
	logger                    *zap.Logger
	event                     *events.Events
}

func NewFinalScoreBoardAdminController(goqu *goqu.Database, logger *zap.Logger, event *events.Events) (*FinalScoreBoardAdminController, error) {
	finalScoreBoardAdminModel, err := models.InitFinalScoreBoardAdminModel(goqu)
	if err != nil {
		return nil, err
	}

	return &FinalScoreBoardAdminController{
		finalScoreBoardAdminModel: &finalScoreBoardAdminModel,
		logger:                    logger,
		event:                     event,
	}, nil

}

// GetScore to send final score after quiz over to admin
// swagger:route GET /v1/final_score/admin FinalScoreForAdmin RequestFinalScoreForAdmin
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
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("user play quiz should be valid string").Error())
	}

	finalScoreBoardData, err := fc.finalScoreBoardAdminModel.GetScoreForAdmin(activeQuizId)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}

	return utils.JSONSuccess(ctx, http.StatusOK, finalScoreBoardData)

}
