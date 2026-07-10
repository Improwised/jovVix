package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type QuizController struct {
	quizModel         *models.QuizModel
	questionModel     *models.QuestionModel
	activeQuizModel   *models.ActiveQuizModel
	quizCategoryModel *models.QuizCategoryModel
	quizSvc           *services.QuizService
	appConfig         *config.AppConfig
	logger            *zap.Logger
}

func InitQuizController(db *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*QuizController, error) {

	quizModel := models.InitQuizModel(db)
	questionModel := models.InitQuestionModel(db, logger)
	activeQuizModel := models.InitActiveQuizModel(db, logger)
	quizCategoryModel := models.InitQuizCategoryModel(db)

	quizSvc := services.NewQuizService(db, logger)

	return &QuizController{
		quizModel:         quizModel,
		questionModel:     questionModel,
		activeQuizModel:   activeQuizModel,
		quizCategoryModel: quizCategoryModel,
		quizSvc:           quizSvc,
		appConfig:         appConfig,
		logger:            logger,
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

func (ctrl *QuizController) CreateQuiz(c *fiber.Ctx) error {
	var quizReq structs.ReqCreateQuiz
	err := json.Unmarshal(c.Body(), &quizReq)
	if err != nil {
		ctrl.logger.Error("validate req error", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(quizReq)
	if err != nil {
		ctrl.logger.Error("validate req error", zap.Any("quizReq", quizReq))
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	// Only configured admin emails may publish public quizzes; silently coerce otherwise.
	isPublic := quizReq.IsPublic
	if isPublic {
		user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
		if !ok || !ctrl.appConfig.Quiz.IsPublicQuizAdmin(user.Email) {
			isPublic = false
		}
	}

	// Category and cover image only exist for the public catalog; drop them
	// whenever the quiz ends up private (same coercion as is_public above).
	categoryId := quizReq.CategoryId
	coverImage := quizReq.CoverImage
	if !isPublic {
		categoryId = ""
		coverImage = ""
	}

	if coverImage != "" {
		if !strings.HasPrefix(coverImage, "data:image/") {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvalidCoverImage)
		}
		if len(coverImage) > constants.MaxCoverImageBytes {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrCoverImageTooLarge)
		}
	}

	if categoryId != "" {
		_, err := ctrl.quizCategoryModel.GetCategoryById(categoryId)
		if err != nil {
			if err == sql.ErrNoRows {
				return utils.JSONFail(c, http.StatusBadRequest, constants.ErrCategoryNotFound)
			}
			ctrl.logger.Error("error fetching category for quiz creation", zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, err.Error())
		}
	}

	quizId, err := ctrl.quizModel.CreateQuiz(quizReq.Title, quizReq.Description, userID, isPublic, categoryId, coverImage)
	if err != nil {
		ctrl.logger.Error("error in creating empty quiz", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrRegisterQuiz)
	}

	return utils.JSONSuccess(c, http.StatusCreated, quizId)
}

// GetPublicQuizzes lists every quiz that has been published publicly.
// Unauthenticated; the homepage uses this to populate "Explore Public Quizzes".
func (ctrl *QuizController) GetPublicQuizzes(c *fiber.Ctx) error {
	quizzes, err := ctrl.quizModel.GetPublicQuizzes()
	if err != nil {
		ctrl.logger.Error("error occured while listing public quizzes", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONSuccess(c, http.StatusOK, quizzes)
}

func (ctrl *QuizController) UpdateQuizSettings(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)

	var quizReq structs.ReqUpdateQuizSettings
	err := json.Unmarshal(c.Body(), &quizReq)
	if err != nil {
		ctrl.logger.Error("validate req error", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(quizReq)
	if err != nil {
		ctrl.logger.Error("validate req error", zap.Any("quizReq", quizReq))
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	err = ctrl.quizSvc.UpdateQuizSettings(quizId, quizReq.Points, quizReq.DurationInSeconds, quizReq.QuestionIds)
	if err != nil {
		ctrl.logger.Error("error in updating quiz settings", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, "error while updating quiz settings")
	}

	return utils.JSONSuccess(c, http.StatusOK, "quiz settings update success")
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

	activeQuiz, err := ctrl.activeQuizModel.GetActiveQuizByQuizIDAndAdminID(quizId, userId)
	if err == nil {
		return utils.JSONSuccess(c, http.StatusAccepted, activeQuiz.ID)
	}
	if err != sql.ErrNoRows {
		ctrl.logger.Error("error checking active demo session", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	quiz, err := ctrl.quizModel.GetQuizById(quizId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuizNotFound)
		}
		ctrl.logger.Error("error fetching quiz for demo session", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	sessionId, err := ctrl.activeQuizModel.CreateActiveQuiz(quiz.Title, quizId, userId, sql.NullTime{}, sql.NullTime{})

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

// GeneratePublicSession lets ANY visitor (guest or registered) host a public quiz.
// Unlike demo_session, it does not require Kratos auth, but it only works for
// quizzes flagged is_public. The starter becomes the session host; whether they
// may also play is decided later at PlayedQuizValidation (creators stay host-only).
func (ctrl *QuizController) GeneratePublicSession(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	quiz, err := ctrl.quizModel.GetQuizById(quizId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuizNotFound)
		}
		ctrl.logger.Error("error fetching quiz for public session", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	if !quiz.IsPublic {
		return utils.JSONFail(c, http.StatusForbidden, constants.ErrQuizNotPublic)
	}

	sessionId, err := ctrl.activeQuizModel.CreateActiveQuiz(quiz.Title, quizId, userId, sql.NullTime{}, sql.NullTime{})
	if err != nil {
		ctrl.logger.Error("error in creating public session", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrCreatingDemoQuiz)
	}

	err = ctrl.activeQuizModel.GetQuestionsCopy(sessionId, quizId)
	if err != nil {
		ctrl.logger.Error("error in creating public session questions", zap.Error(err))
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
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrActiveDeleteQuiz)
	}

	err = ctrl.quizSvc.DeleteQuizById(quizId)
	if err != nil {
		ctrl.logger.Error("error occured while deleting quiz", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	ctrl.logger.Debug("QuizController.DeleteQuizById success", zap.Any(constants.QuizId, quizId))
	return utils.JSONSuccess(c, http.StatusOK, "success")
}
