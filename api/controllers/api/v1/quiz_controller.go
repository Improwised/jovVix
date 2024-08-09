package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/helpers/calculations"
	quiz_helper "github.com/Improwised/quizz-app/api/helpers/quiz"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type QuizController struct {
	helper                  *quiz_helper.HelperGroup
	userPlayedQuizModel     *models.UserPlayedQuizModel
	logger                  *zap.Logger
	event                   *events.Events
	answersSubmittedByUsers chan models.User
}

func InitQuizController(db *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher, helper *quiz_helper.HelperGroup, answersSubmittedByUsers chan models.User) *QuizController {

	userPlayedQuizModel := models.InitUserPlayedQuizModel(db)

	return &QuizController{
		helper:                  helper,
		userPlayedQuizModel:     userPlayedQuizModel,
		logger:                  logger,
		event:                   event,
		answersSubmittedByUsers: answersSubmittedByUsers,
	}
}

func (ctrl *QuizController) GetQuizByUser(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	quizzes, err := ctrl.helper.QuizModel.GetAllQuizzesActivity(userID)

	if err != nil {
		ctrl.logger.Error("error occured while getting all quiz activity", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)

}

func (ctrl *QuizController) CreateQuizSession(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	quizzes, err := ctrl.helper.QuizModel.GetAllQuizzesActivity(userID)

	if err != nil {
		ctrl.logger.Error("error occured while getting all quiz activity", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)
}

func (ctrl *QuizController) SetAnswer(c *fiber.Ctx) error {
	currentQuiz := c.Query(constants.CurrentUserQuiz)
	ctrl.logger.Debug("QuizController.SetAnswer called", zap.Any("currentQuiz", currentQuiz))

	// validations
	if currentQuiz == "" {
		ctrl.logger.Error(constants.ErrQuizNotFound)
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuizNotFound)
	}

	currentQuizId, err := uuid.Parse(currentQuiz)
	if err != nil {
		ctrl.logger.Error("invalid UUID")
		return utils.JSONFail(c, http.StatusBadRequest, "invalid UUID")
	}

	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
	if !ok {
		ctrl.logger.Error("Unable to convert to user-model type from locals")
		return utils.JSONFail(c, http.StatusInternalServerError, "Unable to convert to user-model type from locals")
	}

	var answer structs.ReqAnswerSubmit

	err = json.Unmarshal(c.Body(), &answer)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(answer)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	// check for question is active or not to receive answers
	ctrl.logger.Debug("userPlayedQuizModel.GetCurrentActiveQuestion called", zap.Any("currentQuiz", currentQuiz))
	currentQuestion, err := ctrl.userPlayedQuizModel.GetCurrentActiveQuestion(currentQuizId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctrl.logger.Error("error during answer submit get current active question", zap.Any("answers", answer), zap.Any("current_quiz_id", currentQuizId))
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerSubmit)
		}
		ctrl.logger.Error("error during answer submit", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.UnknownError)
	}
	ctrl.logger.Debug("userPlayedQuizModel.GetCurrentActiveQuestion success", zap.Any("currentQuestion", currentQuestion))

	if currentQuestion != answer.QuestionId {
		ctrl.logger.Error(constants.ErrQuestionNotActive)
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuestionNotActive)
	}

	appConfigs := config.GetConfig()
	db, err := database.Connect(appConfigs.DB)
	if err != nil {
		ctrl.logger.Error("unable to connect with the database", zap.Error(err))
	}

	// calculate points
	points, score, err := calculations.CalculatePointsAndScore(answer, db, ctrl.logger)
	if err != nil {
		if err == sql.ErrNoRows {
			ctrl.logger.Error("error during answer submit", zap.Any("answers", answer), zap.Any("current_quiz_id", currentQuizId))
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerSubmit)
		}

		ctrl.logger.Error("error during answer submit", zap.Error(err))

		return utils.JSONFail(c, http.StatusBadRequest, constants.UnknownError)
	}

	// core logic
	err = ctrl.helper.UserQuizResponseModel.SubmitAnswer(currentQuizId, answer, points, score)

	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerAlreadySubmitted)
		}
		ctrl.logger.Error("error during answer submit", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.UnknownError)
	}

	ctrl.answersSubmittedByUsers <- user

	return utils.JSONSuccess(c, http.StatusAccepted, nil)
}

// GetAdminUploadedQuizzes for getting quiz details uploaded by Admin
// swagger:route GET /v1/admin/quizzes/list AdminUploadedQuizzes none
//
// Get details of quizzes uploaded by Admin.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseAdminUploadedQuizDetails
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
func (ctrl *QuizController) GetAdminUploadedQuizzes(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	quizzes, err := ctrl.helper.QuizModel.GetQuizzesByAdmin(userID)

	if err != nil {
		ctrl.logger.Error("error occured while getting quizzes by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)
}

func (ctrl *QuizController) Terminate(c *fiber.Ctx) error {
	return utils.JSONSuccess(c, http.StatusOK, nil)
}
