package v1

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type QuizController struct {
	userPlayedQuizModel   *models.UserPlayedQuizModel
	questionModel         *models.QuestionModel
	userQuizResponseModel *models.UserQuizResponseModel
	quizModel             *models.QuizModel
	activeQuizModel       *models.ActiveQuizModel
	presignedURLSvc       *services.PresignURLService
	logger                *zap.Logger
	event                 *events.Events
}

func InitQuizController(db *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher, appConfig *config.AppConfig) (*QuizController, error) {

	userPlayedQuizModel := models.InitUserPlayedQuizModel(db)
	questionModel := models.InitQuestionModel(db, logger)
	quizModel := models.InitQuizModel(db)
	userQuizResponseModel := models.InitUserQuizResponseModel(db)
	activeQuizModel := models.InitActiveQuizModel(db, logger)

	presignedURLSvc, err := services.NewFileUploadServices(appConfig.AWS.BucketName)
	if err != nil {
		return nil, err
	}

	return &QuizController{
		userPlayedQuizModel:   userPlayedQuizModel,
		questionModel:         questionModel,
		userQuizResponseModel: userQuizResponseModel,
		quizModel:             quizModel,
		activeQuizModel:       activeQuizModel,
		presignedURLSvc:       presignedURLSvc,
		logger:                logger,
		event:                 event,
	}, nil
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

	services.ProcessAnalyticsData(quizAnalysis, qc.presignedURLSvc, qc.logger)

	return utils.JSONSuccess(c, http.StatusOK, quizAnalysis)
}

func (qc *QuizController) ListQuizzesAnalysis(c *fiber.Ctx) error {

	type resQuizAnalysisList struct {
		Data  []models.QuizzesAnalysis
		Count int64
	}

	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	filters := c.Queries()
	var page int
	var err error

	if val, isSet := filters[constants.PageNumberQueryParam]; isSet {
		page, err = strconv.Atoi(val)
		if err != nil {
			return utils.JSONFail(c, http.StatusBadRequest, "Enter page number in integer only.")
		}

		if page <= 0 {
			page = 1
		}
	} else {
		page = 1
	}

	quizzes, count, err := qc.quizModel.ListQuizzesAnalysis(filters[constants.NameQueryParam], filters[constants.OrderQueryParam], filters[constants.OrderByQueryParam], filters["date"], userID, page)

	if err != nil {
		qc.logger.Error("error occured while listing quizzes for analysis", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, resQuizAnalysisList{Data: quizzes, Count: count})
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

func (ctrl *QuizController) ListQuestionByQuizId(c *fiber.Ctx) error {
	QuizId := c.Params(constants.QuizId)
	Query := c.Queries()
	ctrl.logger.Debug("QuizController.ListQuestionByQuizId called", zap.Any(constants.QuizId, QuizId), zap.Any("Query", Query))

	questions, err := ctrl.questionModel.ListQuestionByQuizId(QuizId, Query[constants.MediaQuery])
	if err != nil {
		ctrl.logger.Error("error occured while getting quizzes by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	ctrl.logger.Debug("QuizController.ListQuestionByQuizId success", zap.Any("questions", questions))
	return utils.JSONSuccess(c, http.StatusOK, questions)
}
