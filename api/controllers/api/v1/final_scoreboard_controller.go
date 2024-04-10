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
// swagger:route GET /v1/finalScore FinalScore RequestFinalScore
//
// Get a finalScore details.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseFinalScore
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (fc *FinalScoreBoardController) GetScore(ctx *fiber.Ctx) error {

	var userPlayedQuiz = ctx.Query(constants.UserPlayedQuiz)

	if !(userPlayedQuiz != "" && len(userPlayedQuiz) == 36) {
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("user play quiz should be valid string").Error())
	}

	finalScoreBoardData, err := fc.finalScoreBoardModel.GetScore(userPlayedQuiz)
	if err != nil {
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}

	return utils.JSONSuccess(ctx, http.StatusOK, finalScoreBoardData)

}
