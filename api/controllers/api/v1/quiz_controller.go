package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Improwised/quizz-app/api/constants"
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
	userPlayedQuizModel     *models.UserPlayedQuizModel
	questionModel           *models.QuestionModel
	userQuizResponseModel   *models.UserQuizResponseModel
	quizModel               *models.QuizModel
	activeQuizModel         *models.ActiveQuizModel
	logger                  *zap.Logger
	event                   *events.Events
	answersSubmittedByUsers chan models.User
}

func InitQuizController(db *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher, answersSubmittedByUsers chan models.User) *QuizController {

	userPlayedQuizModel := models.InitUserPlayedQuizModel(db)
	questionModel := models.InitQuestionModel(db, logger)
	quizModel := models.InitQuizModel(db)
	userQuizResponseModel := models.InitUserQuizResponseModel(db)
	activeQuizModel := models.InitActiveQuizModel(db, logger)

	return &QuizController{
		userPlayedQuizModel:     userPlayedQuizModel,
		questionModel:           questionModel,
		userQuizResponseModel:   userQuizResponseModel,
		quizModel:               quizModel,
		activeQuizModel:         activeQuizModel,
		logger:                  logger,
		event:                   event,
		answersSubmittedByUsers: answersSubmittedByUsers,
	}
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

	answers, answerPoints, answerDurationInSeconds, questionType, err := ctrl.questionModel.GetAnswersPointsDurationType(answer.QuestionId.String())
	if err != nil {
		ctrl.logger.Error("error while get answer, points, duration and type")
		return utils.JSONFail(c, http.StatusBadRequest, "error while get answer, points, duration and type")
	}

	// calculate points
	points, score := utils.CalculatePointsAndScore(answer, answers, answerPoints, answerDurationInSeconds, questionType)

	// submit answer
	err = ctrl.userQuizResponseModel.SubmitAnswer(currentQuizId, answer, points, score)
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

	quizzes, err := ctrl.quizModel.GetQuizzesByAdmin(userID)

	if err != nil {
		ctrl.logger.Error("error occured while getting quizzes by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)
}

func (ctrl *QuizController) Terminate(c *fiber.Ctx) error {
	return utils.JSONSuccess(c, http.StatusOK, nil)
}

func (qc *QuizController) GetQuizAnalysis(c *fiber.Ctx) error {

	activeQuizId := c.Params(constants.ActiveQuizId)

	quizAnalysis, err := qc.quizModel.GetQuizAnalysis(activeQuizId)

	if err != nil {
		qc.logger.Error("error while get quiz analysis", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizAnalysis)
}

func (qc *QuizController) ListQuizzesAnalysis(c *fiber.Ctx) error {

	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	quizzes, err := qc.quizModel.ListQuizzesAnalysis(userID)

	if err != nil {
		qc.logger.Error("error occured while listing quizzes for analysis", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)
}

func (ctrl *QuizController) CreateQuizByCsv(c *fiber.Ctx) error {

	quizTitle := c.Params(constants.QuizTitle)
	quizDescription := c.FormValue("description")

	if quizTitle == "" {
		ctrl.logger.Error("quiz-title not found")
		return utils.JSONSuccess(c, http.StatusBadRequest, constants.QuizTitleRequired)
	}

	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	filePath := quizUtilsHelper.GetString(c.Locals(constants.FileName))

	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			ctrl.logger.Error("error in deleting file", zap.Error(err))
			return
		}
	}()

	questions, err := utils.ValidateCSVFileFormat(filePath)
	if err != nil {
		ctrl.logger.Error("file validation failed", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validQuestions, err := utils.ExtractQuestionsFromCSV(questions)
	if err != nil {
		ctrl.logger.Error("file validation failed", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrParsingFile)
	}

	quizId, err := ctrl.questionModel.RegisterQuestions(userID, quizTitle, quizDescription, validQuestions)
	if err != nil {
		ctrl.logger.Error("error in creating quiz", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrRegisterQuiz)
	}

	return utils.JSONSuccess(c, http.StatusAccepted, quizId)
}

func (ctrl *QuizController) GenerateDemoSession(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	sessionId, err := ctrl.activeQuizModel.CreateActiveQuiz("demo session", quizId, userId, sql.NullTime{}, sql.NullTime{})

	if err != nil {
		ctrl.logger.Error("error in creating demo session", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrCreatingDemoQuiz)
	}

	err = ctrl.activeQuizModel.GetQuestionsCopy(sessionId, quizId)
	if err != nil {
		ctrl.logger.Error("error in creating demo session questions", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrCreatingDemoQuiz)
	}

	return utils.JSONSuccess(c, http.StatusAccepted, sessionId)
}
