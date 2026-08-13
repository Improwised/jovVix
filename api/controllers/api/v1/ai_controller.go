package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type AIController struct {
	questionModel *models.QuestionModel
	quizSvc       *services.QuizService
	aiSvc         *services.AIQuizService
	appConfig     *config.AppConfig
	logger        *zap.Logger
}

func InitAIController(db *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*AIController, error) {
	questionModel := models.InitQuestionModel(db, logger)
	aiSvc := services.NewAIQuizService(logger, &appConfig.AI)

	return &AIController{
		questionModel: questionModel,
		quizSvc:       services.NewQuizService(db, logger),
		aiSvc:         aiSvc,
		appConfig:     appConfig,
		logger:        logger,
	}, nil
}

func (ctrl *AIController) resolveCredentials(c *fiber.Ctx) (services.AICredentials, error) {
	return services.ResolveAICredentials(
		c.Get(constants.HeaderAIBaseUrl),
		c.Get(constants.HeaderAIApiKey),
		c.Get(constants.HeaderAIModel),
	)
}

func (ctrl *AIController) credentialError(c *fiber.Ctx, err error) error {
	if err.Error() == constants.ErrAINotConfigured {
		ctrl.logger.Error("ai request received while unconfigured")
		return utils.JSONError(c, http.StatusServiceUnavailable, err.Error())
	}
	return utils.JSONFail(c, http.StatusBadRequest, err.Error())
}

// existingQuestions feeds the prompt's avoid list. A lookup failure is logged and
// treated as an empty quiz: generating without the list beats not generating.
func (ctrl *AIController) existingQuestions(quizId string) []string {
	if quizId == "" {
		return nil
	}

	questions, err := ctrl.questionModel.ListQuestionTextsByQuizId(quizId)
	if err != nil {
		ctrl.logger.Error("could not read existing questions for the ai avoid list",
			zap.String("quiz_id", quizId),
			zap.Error(err))
		return nil
	}

	return questions
}

func (ctrl *AIController) resolveQuestionOptions() utils.AIQuestionOptions {
	duration, err := strconv.Atoi(strings.TrimSpace(ctrl.appConfig.Quiz.QuestionTimeLimit))
	if err != nil || duration <= 0 {
		duration = constants.AIDefaultQuestionDuration
	}

	points := ctrl.appConfig.Quiz.DefaultQuestionPoints
	if points <= constants.MinimumPoints {
		points = constants.AIDefaultQuestionPoints
	}

	return utils.AIQuestionOptions{
		Duration:     duration,
		Points:       points,
		MaxQuestions: ctrl.appConfig.AI.QuestionLimit(),
	}
}

// swagger:route GET /v1/ai/status AI RequestAIStatus
//
// Report whether the server has a provider configured.
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseAIStatus
//	     401: GenericResFailConflict
func (ctrl *AIController) Status(c *fiber.Ctx) error {
	return utils.JSONSuccess(c, http.StatusOK, structs.ResAIStatus{
		MaxQuestions: ctrl.appConfig.AI.QuestionLimit(),
	})
}

// swagger:route POST /v1/ai/test AI RequestAITestConnection
//
// Test an AI provider connection.
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseAITestConnection
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  502: GenericResError
//			  503: GenericResError
func (ctrl *AIController) TestConnection(c *fiber.Ctx) error {
	cred, err := ctrl.resolveCredentials(c)
	if err != nil {
		return ctrl.credentialError(c, err)
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), constants.AITestTimeoutSeconds*time.Second)
	defer cancel()

	elapsed, err := ctrl.aiSvc.TestConnection(ctx, cred)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadGateway, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, structs.ResAITestConnection{
		Ok:        true,
		Model:     cred.Model,
		Endpoint:  cred.CompletionsURL(),
		LatencyMs: int(elapsed.Milliseconds()),
	})
}

// swagger:route GET /v1/ai/models AI RequestAIModels
//
// List the models available from an AI provider.
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseAIModels
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  502: GenericResError
//			  503: GenericResError
func (ctrl *AIController) ListModels(c *fiber.Ctx) error {
	cred, err := services.ResolveAICredentialsForListing(
		c.Get(constants.HeaderAIBaseUrl),
		c.Get(constants.HeaderAIApiKey),
	)
	if err != nil {
		return ctrl.credentialError(c, err)
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), constants.AITestTimeoutSeconds*time.Second)
	defer cancel()

	models, err := ctrl.aiSvc.ListModels(ctx, cred)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadGateway, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, structs.ResAIModels{
		Models: models,
	})
}

// swagger:route POST /v1/ai/questions/generate AI RequestGenerateAIQuestions
//
// Generate quiz questions for a topic with AI.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseGenerateAIQuestions
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
func (ctrl *AIController) GenerateQuestions(c *fiber.Ctx) error {
	return ctrl.generate(c, "")
} //

// swagger:route POST /v1/quizzes/{quiz_id}/questions/ai/generate AI RequestGenerateAIQuestionsForQuiz
//
// Generate quiz questions with AI for an existing quiz.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseGenerateAIQuestions
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
func (ctrl *AIController) GenerateQuestionsForQuiz(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)
	if quizId == "" {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAIMissingQuizId)
	}
	return ctrl.generate(c, quizId)
}

// generate serves both entry points. A quiz id means the questions are destined for
// an existing quiz, so the ones it already holds are sent along as an avoid list.
func (ctrl *AIController) generate(c *fiber.Ctx, quizId string) error {
	cred, err := ctrl.resolveCredentials(c)
	if err != nil {
		return ctrl.credentialError(c, err)
	}

	var req structs.ReqGenerateAIQuestions
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		ctrl.logger.Error("validate req error", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	req.Topic = strings.TrimSpace(req.Topic)
	req.Difficulty = strings.ToLower(strings.TrimSpace(req.Difficulty))
	req.Language = strings.ToLower(strings.TrimSpace(req.Language))
	if req.Language == "" {
		req.Language = constants.AIDefaultLanguage
	}

	if strings.ContainsAny(req.Topic, "\r\n") {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAITopicSingleLine)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		ctrl.logger.Error("validate req error", zap.String("validation", utils.ValidatorErrorString(err)))
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	limit := ctrl.appConfig.AI.QuestionLimit()
	if req.NumberOfQuestions > limit {
		return utils.JSONFail(c, http.StatusBadRequest, fmt.Sprintf(constants.ErrAITooManyQuestions, limit))
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), ctrl.appConfig.AI.Timeout())
	defer cancel()

	generation, err := ctrl.aiSvc.GenerateQuestions(ctx, cred, req, ctrl.existingQuestions(quizId))
	if err != nil {
		switch err.Error() {
		case constants.ErrAINotConfigured:
			return utils.JSONError(c, http.StatusServiceUnavailable, err.Error())
		case constants.ErrAIRateLimited, constants.ErrAIBudgetExhausted:
			return utils.JSONError(c, http.StatusTooManyRequests, err.Error())
		default:
			return utils.JSONError(c, http.StatusBadGateway, err.Error())
		}
	}

	result, err := utils.NormalizeAIQuestions(generation.Questions, ctrl.resolveQuestionOptions())
	if err != nil {
		ctrl.logger.Error("no usable questions in ai response",
			zap.Strings("issues", result.Issues),
			zap.Error(err))
		return utils.JSONError(c, http.StatusBadGateway, constants.ErrAINoValidQuestions)
	}

	if len(result.Issues) > 0 {
		ctrl.logger.Warn("some generated questions were dropped", zap.Strings("issues", result.Issues))
	}

	var notices []string
	if len(result.Accepted) < req.NumberOfQuestions {
		notices = append(notices, fmt.Sprintf(constants.NoticeAIFewerQuestions, len(result.Accepted), req.NumberOfQuestions))
	}
	if len(result.Issues) > 0 {
		notices = append(notices, fmt.Sprintf(constants.NoticeAIDroppedQuestions, len(result.Issues)))
	}

	return utils.JSONSuccess(c, http.StatusOK, structs.ResGenerateAIQuestions{
		Topic:                req.Topic,
		Difficulty:           req.Difficulty,
		Language:             req.Language,
		Generated:            len(result.Accepted),
		SuggestedTitle:       utils.SanitizeAITitle(generation.Title, req.Topic),
		SuggestedDescription: utils.SanitizeAIDescription(generation.Description),
		Notice:               strings.Join(notices, " "),
		Questions:            result.Accepted,
	})
}

// swagger:route POST /v1/ai/quizzes AI RequestCreateQuizFromAI
//
// Create a quiz from AI generated questions.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  201: ResponseQuizCreated
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
func (ctrl *AIController) CreateQuizFromAI(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	if userID == "" {
		return utils.JSONFail(c, http.StatusUnauthorized, constants.ErrUnauthenticated)
	}

	var req structs.ReqCreateQuizFromAI
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		ctrl.logger.Error("validate req error", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		ctrl.logger.Error("validate req error", zap.Any("title", req.Title))
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	result, err := utils.NormalizeAIQuestions(req.Questions, ctrl.resolveQuestionOptions())
	if err != nil || len(result.Issues) > 0 || len(result.Questions) != len(req.Questions) {
		ctrl.logger.Error("rejected ai questions on create",
			zap.Strings("issues", result.Issues),
			zap.Int("submitted", len(req.Questions)),
			zap.Int("accepted", len(result.Questions)),
			zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAIInvalidQuestions)
	}

	quizId, err := ctrl.questionModel.RegisterQuizAndQuestions(userID, req.Title, req.Description, result.Questions)
	if err != nil {
		ctrl.logger.Error("error in creating quiz from ai questions", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrRegisterQuiz)
	}

	return utils.JSONSuccess(c, http.StatusCreated, quizId)
}

// swagger:route POST /v1/quizzes/{quiz_id}/questions/ai AI RequestAppendAIQuestions
//
// Append AI generated questions to an existing quiz.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  201: ResponseAppendAIQuestions
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
func (ctrl *AIController) AppendAIQuestions(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)
	if quizId == "" {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAIMissingQuizId)
	}

	var req structs.ReqAppendAIQuestions
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		ctrl.logger.Error("validate req error", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		ctrl.logger.Error("validate req error", zap.Int("questions", len(req.Questions)))
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	options := ctrl.resolveQuestionOptions()
	if req.DurationInSeconds > 0 {
		options.Duration = req.DurationInSeconds
	}
	if req.Points > constants.MinimumPoints {
		options.Points = req.Points
	}

	result, err := utils.NormalizeAIQuestions(req.Questions, options)
	if err != nil || len(result.Issues) > 0 || len(result.Questions) != len(req.Questions) {
		ctrl.logger.Error("rejected ai questions on append",
			zap.String("quiz_id", quizId),
			zap.Strings("issues", result.Issues),
			zap.Int("submitted", len(req.Questions)),
			zap.Int("accepted", len(result.Questions)),
			zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAIInvalidQuestions)
	}

	questionIds, err := ctrl.quizSvc.AppendQuestionsToQuiz(quizId, result.Questions)
	if err != nil {
		ctrl.logger.Error("error in appending ai questions to quiz",
			zap.String("quiz_id", quizId),
			zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrAIAppendFailed)
	}

	if len(questionIds) == 0 {
		ctrl.logger.Error("append returned no question ids", zap.String("quiz_id", quizId))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrAIAppendFailed)
	}

	return utils.JSONSuccess(c, http.StatusCreated, structs.ResAppendAIQuestions{
		Added:       len(questionIds),
		QuestionIds: questionIds,
	})
}
