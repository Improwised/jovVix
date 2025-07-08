package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/events"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/pkg/watermill"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type QuestionController struct {
	questionModel   *models.QuestionModel
	activeQuizModel *models.ActiveQuizModel
	presignedURLSvc *services.PresignURLService
	quizSvc         *services.QuizService
	logger          *zap.Logger
	event           *events.Events
}

func InitQuestionController(db *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher, appConfig *config.AppConfig) (*QuestionController, error) {

	questionModel := models.InitQuestionModel(db, logger)
	activeQuizModel := models.InitActiveQuizModel(db, logger)

	presignedURLSvc, err := services.NewFileUploadServices(&appConfig.AWS)
	if err != nil {
		return nil, err
	}

	quizSvc := services.NewQuizService(db, logger)

	return &QuestionController{
		questionModel:   questionModel,
		activeQuizModel: activeQuizModel,
		presignedURLSvc: presignedURLSvc,
		quizSvc:         quizSvc,
		logger:          logger,
		event:           event,
	}, nil
}

// ListQuestionByQuizId to list all questions of quiz with `is_active_quiz_present` and `quiz_played_count`.
// swagger:route GET /v1/quizzes/{quiz_id}/questions Question RequestListQuestionByQuizId
//
// List all questions of quiz with `is_active_quiz_present` and `quiz_played_count`.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseListQuestionByQuizId
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (ctrl *QuestionController) ListQuestionsWithAnswerByQuizId(c *fiber.Ctx) error {
	QuizId := c.Params(constants.QuizId)
	Query := c.Queries()
	permission := c.Locals(constants.ContextQuizPermission).(string)
	ctrl.logger.Debug("QuestionController.ListQuestionsWithAnswerByQuizId called", zap.Any(constants.QuizId, QuizId), zap.Any("Query", Query))

	isActiveQuizPresent, err := ctrl.activeQuizModel.IsActiveQuizPresent(QuizId)
	if err != nil {
		ctrl.logger.Error("error occured while getting questions by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	questions, quizPlayedcount, err := ctrl.questionModel.ListQuestionsWithAnswerByQuizId(QuizId, Query[constants.MediaQuery])
	if err != nil {
		ctrl.logger.Error("error occured while getting questions by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	services.ProcessAnalyticsData(questions, ctrl.presignedURLSvc, ctrl.logger)

	ctrl.logger.Debug("QuestionController.ListQuestionsWithAnswerByQuizId success", zap.Any("questions", structs.ResQuestionAnalytics{Data: questions, QuizPlayedCount: quizPlayedcount}), zap.Any("quizPlayedcount", quizPlayedcount))
	return utils.JSONSuccess(c, http.StatusOK, structs.ResQuestionAnalytics{Data: questions, QuizPlayedCount: quizPlayedcount, IsQuizEditable: !isActiveQuizPresent, Permission: permission})
}

// GetQuestionById to get question and thier options with answer.
// swagger:route GET /v1/quizzes/{quiz_id}/questions/{question_id} Question RequestGetQuestionById
//
// Get question and thier options with answer.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseGetQuestionById
//	     401: GenericResFailConflict
//		  500: GenericResError
func (ctrl *QuestionController) GetQuestionById(c *fiber.Ctx) error {
	QuestionId := c.Params(constants.QuestionId)
	ctrl.logger.Debug("QuestionController.GetQuestionById called", zap.Any(constants.QuestionId, QuestionId))

	question, err := ctrl.questionModel.GetQuestionById(QuestionId)
	if err != nil {
		ctrl.logger.Error("error occured while getting question by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	if question.QuestionsMedia == "image" {
		presignedURL, err := ctrl.presignedURLSvc.GetPresignedURL(question.Resource, 5*time.Minute)
		if err != nil {
			ctrl.logger.Error("error while getting presigned url", zap.Error(err))
		}
		question.Resource = presignedURL
	}

	if question.OptionsMedia == "image" {
		for key, value := range question.Options {
			presignedURL, err := ctrl.presignedURLSvc.GetPresignedURL(value, 5*time.Minute)
			if err != nil {
				ctrl.logger.Error("error while getting presigned url", zap.Error(err))
			}
			question.Options[key] = presignedURL
		}
	}

	ctrl.logger.Debug("QuestionController.GetQuestionById success", zap.Any("question", question))
	return utils.JSONSuccess(c, http.StatusOK, question)
}

// UpdateQuestionById to update question and thier options with answer.
// swagger:route PUT /v1/quizzes/{quiz_id}/questions/{question_id} Question RequestUpdateQuestionById
//
// Update question and thier options with answer.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseOkWithMessage
//	     401: GenericResFailConflict
//		  500: GenericResError
//	     400: GenericResFailNotFound
func (ctrl *QuestionController) UpdateQuestionById(c *fiber.Ctx) error {
	QuestionId := c.Params(constants.QuestionId)
	ctrl.logger.Debug("QuizController.UpdateQuestionById called", zap.Any(constants.QuestionId, QuestionId))

	ctrl.logger.Debug("validate req", zap.Any("Body", c.Body()))
	var questionReq structs.ReqUpdateQuestion
	err := json.Unmarshal(c.Body(), &questionReq)
	if err != nil {
		ctrl.logger.Error("validate req error", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(questionReq)
	if err != nil {
		ctrl.logger.Error("validate req error", zap.Any("questionReq", questionReq))
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	err = ctrl.questionModel.UpdateQuestionById(QuestionId, models.Question{
		Question:          questionReq.Question,
		Type:              questionReq.Type,
		Options:           questionReq.Options,
		Answers:           questionReq.Answers,
		Points:            questionReq.Points,
		DurationInSeconds: questionReq.DurationInSeconds,
		QuestionMedia:     questionReq.QuestionMedia,
		OptionsMedia:      questionReq.OptionsMedia,
		Resource:          sql.NullString{String: questionReq.Resource, Valid: true},
	})
	if err != nil {
		ctrl.logger.Error("error occured while update question by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	ctrl.logger.Debug("QuizController.UpdateQuestionById success", zap.Any("QuestionId", QuestionId))
	return utils.JSONSuccess(c, http.StatusOK, "question update success")
}

// DeleteQuestionById to delete question only if no active quiz is present.
// swagger:route DELETE /v1/quizzes/{quiz_id}/questions/{question_id} Question RequestDeleteQuestionById
//
// Delete question only if no active quiz is present.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseOkWithMessage
//	     401: GenericResFailConflict
//		  500: GenericResError
//	     400: GenericResFailNotFound
func (ctrl *QuestionController) DeleteQuestionById(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)
	ctrl.logger.Debug("QuizController.DeleteQuizById called", zap.Any(constants.QuizId, quizId))

	questionId := c.Params(constants.QuestionId)
	ctrl.logger.Debug("QuizController.DeleteQuizById called", zap.Any(constants.QuestionId, questionId))

	isActiveQuizPresent, err := ctrl.activeQuizModel.IsActiveQuizPresent(quizId)
	if err != nil {
		ctrl.logger.Error("error occured while getting questions by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}
	if isActiveQuizPresent {
		ctrl.logger.Error("error occured while getting questions by admin", zap.Error(err))
		return utils.JSONError(c, http.StatusBadRequest, constants.InvalidCredentials)
	}

	err = ctrl.quizSvc.DeleteQuestionById(questionId)
	if err != nil {
		ctrl.logger.Error("error occured while deleting quiz", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	ctrl.logger.Debug("QuizController.ListQuestionsWithAnswerByQuizId success", zap.Any(constants.QuizId, quizId))
	return utils.JSONSuccess(c, http.StatusOK, "success")
}
