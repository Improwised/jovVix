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

type FinalScoreBoardController struct {
	finalScoreBoardModel *models.FinalScoreBoardModel
	logger               *zap.Logger
	event                *events.Events
}

func NewFinalScoreBoardController(goqu *goqu.Database, logger *zap.Logger, event *events.Events) (*FinalScoreBoardController, error) {
	finalScoreBoardModel, err := models.InitFinalScoreBoardModel(goqu)
	if err != nil {
		return nil, err
	}

	return &FinalScoreBoardController{
		finalScoreBoardModel: &finalScoreBoardModel,
		logger:               logger,
		event:                event,
	}, nil

}

// GetScore to send final score after quiz over
// swagger:route GET /v1/final_score/user FinalScore RequestFinalScoreForUser
//
// Get a finalScore details for user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseFinalScoreForUser
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (fc *FinalScoreBoardController) GetScore(ctx *fiber.Ctx) error {
	userPlayedQuiz := ctx.Query(constants.UserPlayedQuiz)
	fc.logger.Debug("FinalScoreBoardController.GetScore called", zap.Any("userPlayedQuiz", userPlayedQuiz))

	if !(userPlayedQuiz != "" && len(userPlayedQuiz) == 36) {
		fc.logger.Debug("user played quiz id is not valid - either empty string or it is not 36 characters long")
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("user play quiz should be valid string").Error())
	}

	fc.logger.Debug("finalScoreBoardModel.GetScore called", zap.Any("userPlayedQuiz", userPlayedQuiz))
	finalScoreBoardData, err := fc.finalScoreBoardModel.GetScore(userPlayedQuiz)
	if err != nil {
		fc.logger.Error("Error while getting final scoreboard for user", zap.Error(err))
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}
	fc.logger.Debug("finalScoreBoardModel.GetScore success", zap.Any("finalScoreBoardData", finalScoreBoardData))
	fc.logger.Debug("FinalScoreBoardController.GetScore success", zap.Any("finalScoreBoardData", finalScoreBoardData))

	return utils.JSONSuccess(ctx, http.StatusOK, finalScoreBoardData)
}
