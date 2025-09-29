package v1

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type QuizController struct {
	quizModel       *models.QuizModel
	questionModel   *models.QuestionModel
	activeQuizModel *models.ActiveQuizModel
	presignedURLSvc *services.PresignURLService
	quizSvc         *services.QuizService
	appConfig       *config.AppConfig
	logger          *zap.Logger
}

func InitQuizController(db *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*QuizController, error) {

	quizModel := models.InitQuizModel(db)
	questionModel := models.InitQuestionModel(db, logger)
	activeQuizModel := models.InitActiveQuizModel(db, logger)

	presignedURLSvc, err := services.NewFileUploadServices(&appConfig.AWS)
	if err != nil {
		return nil, err
	}

	quizSvc := services.NewQuizService(db, logger)

	return &QuizController{
		quizModel:       quizModel,
		questionModel:   questionModel,
		activeQuizModel: activeQuizModel,
		presignedURLSvc: presignedURLSvc,
		quizSvc:         quizSvc,
		appConfig:       appConfig,
		logger:          logger,
	}, nil
}

// GetAdminUploadedQuizzes for getting quiz details uploaded by Admin
// swagger:route GET /v1/quizzes Quiz GetAdminUploadedQuizzes
//
// Get details of quizzes uploaded by Admin.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseAdminUploadedQuiz
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

// GetQuizAnalysis for getting quiz details hosted by Admin
// swagger:route GET /v1/admin/reports/{active_quiz_id}/analysis Reports RequestGetQuizAnalysis
//
// Get details of quizzes host by Admin.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseGetQuizAnalysis
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
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

// GetQuizAnalysis for getting quiz list hosted by Admin
// swagger:route GET /v1/admin/reports/list Reports RequestListQuizzesAnalysis
//
// Get details of quizzes list host by Admin.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: RsponseListQuizzesAnalysis
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
func (qc *QuizController) ListQuizzesAnalysis(c *fiber.Ctx) error {

	type resQuizAnalysisList struct {
		Data  []models.QuizzesAnalysis `json:"data"`
		Count int64                    `json:"count"`
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

// CreateQuizByCsv a new quiz by uploading a CSV file
// swagger:route POST /v1/quizzes/{quiz_title}/upload Quiz RequestQuizCreated
//
// Create a new quiz by uploading a CSV file.
//
//			Consumes:
//			- multipart/form-data
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseQuizCreated
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
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

	validQuestions, err := utils.ExtractQuestionsFromCSV(questions, ctrl.appConfig.Quiz.QuestionTimeLimit)
	if err != nil {
		ctrl.logger.Error("file validation failed", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrParsingFile)
	}

	quizId, err := ctrl.questionModel.RegisterQuizAndQuestions(userID, quizTitle, quizDescription, validQuestions)
	if err != nil {
		ctrl.logger.Error("error in creating quiz", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrRegisterQuiz)
	}

	return utils.JSONSuccess(c, http.StatusAccepted, quizId)
}

// GenerateDemoSession to create quiz active for user.
// swagger:route POST /v1/quizzes/{quiz_id}/demo_session Quiz RequestGenerateDemoSession
//
// Create quiz active for user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseGenerateDemoSession
//	     400: GenericResFailNotFound
//		  500: GenericResError
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

// DeleteQuizById to delete quiz that created by user (if no active quiz is present).
// swagger:route DELETE /v1/quizzes/{quiz_id} Quiz DeleteQuizById
//
// Delete quiz that created by user (if no active quiz is present).
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseOkWithMessage
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (ctrl *QuizController) DeleteQuizById(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)
	ctrl.logger.Debug("QuizController.DeleteQuizById called", zap.Any(constants.QuizId, quizId))

	isActiveQuizPresent, err := ctrl.activeQuizModel.IsActiveQuizPresent(quizId)
	if err != nil {
		ctrl.logger.Error("error occured while getting is active quiz is present or not", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	if isActiveQuizPresent {
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidCredentials)
	}

	err = ctrl.quizSvc.DeleteQuizById(quizId)
	if err != nil {
		ctrl.logger.Error("error occured while deleting quiz", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	ctrl.logger.Debug("QuizController.DeleteQuizById success", zap.Any(constants.QuizId, quizId))
	return utils.JSONSuccess(c, http.StatusOK, "success")
}

func (qc *QuizController) DownloadReport(c *fiber.Ctx) error {
	questionType := c.Query(constants.QuestionType)
	activeQuizId := c.Params(constants.ActiveQuizId)
	contentType := c.Get(constants.ContentType)

	quizAnalysis, err := qc.quizModel.GetQuizAnalysis(activeQuizId)

	if err != nil {
		qc.logger.Error("error while get quiz analysis", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	if contentType == constants.CsvAcceptHeader {

		filePath, err := services.GenerateCsv(qc.appConfig.FilePath, activeQuizId, quizAnalysis, func(a []models.QuizAnalysis) [][]string {
			return utils.GetQuestionsData(a, questionType)
		})

		if err != nil {
			return err
		}
		return utils.CsvFileResponse(c, filePath, activeQuizId+".csv")
	}

	return utils.JSONSuccess(c, http.StatusOK, quizAnalysis)
}
