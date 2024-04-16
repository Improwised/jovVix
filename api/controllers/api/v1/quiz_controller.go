package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	quiz_helper "github.com/Improwised/quizz-app/api/helpers/quiz"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type QuizController struct {
	helper *quiz_helper.HelperGroup
	logger *zap.Logger
	event  *events.Events
}

func InitQuizController(
	logger *zap.Logger,
	event *events.Events,
	pub *watermill.WatermillPublisher,
	helper *quiz_helper.HelperGroup,
) *QuizController {
	return &QuizController{
		helper: helper,
		logger: logger,
		event:  event,
	}
}

func (ctrl *QuizController) GetQuizByUser(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	quizzes, err := ctrl.helper.QuizModel.GetAllQuizzesActivity(userID)

	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)

}

func (ctrl *QuizController) CreateQuizSession(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	quizzes, err := ctrl.helper.QuizModel.GetAllQuizzesActivity(userID)

	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)
}

func (ctrl *QuizController) SetAnswer(c *fiber.Ctx) error {
	currentQuiz := c.Cookies(constants.CurrentUserQuiz)

	// validations
	if currentQuiz == "" {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuizNotFound)
	}
	currentQuizId := uuid.MustParse(currentQuiz)

	var answer structs.ReqAnswerSubmit

	err := json.Unmarshal(c.Body(), &answer)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(answer)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	// check for question is active or not to receive answers
	currentQuestion, err := ctrl.helper.UserPlayedQuizModel.GetCurrentActiveQuestion(currentQuizId)

	if err != nil {
		if err == sql.ErrNoRows {
			ctrl.logger.Error("error during answer submit get current active question", zap.Any("answers", answer), zap.Any("current_quiz_id", currentQuizId))
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerSubmit)
		}

		ctrl.logger.Error("error during answer submit", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.UnknownError)
	}

	if currentQuestion != answer.QuestionId {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuestionNotActive)
	}

	// calculate score
	score, err := ctrl.helper.QuestionModel.CalculateScore(answer)

	if err != nil {
		if err == sql.ErrNoRows {
			ctrl.logger.Error("error during answer submit", zap.Any("answers", answer), zap.Any("current_quiz_id", currentQuizId))
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerSubmit)
		}

		ctrl.logger.Error("error during answer submit", zap.Error(err))

		return utils.JSONFail(c, http.StatusBadRequest, constants.UnknownError)
	}

	// core logic
	err = ctrl.helper.UserQuizResponseModel.SubmitAnswer(currentQuizId, answer, score)

	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerAlreadySubmitted)
		}
		ctrl.logger.Error("error during answer submit", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.UnknownError)
	}

	return utils.JSONSuccess(c, http.StatusAccepted, nil)
}

func (ctrl *QuizController) Terminate(c *fiber.Ctx) error {
	return utils.JSONSuccess(c, http.StatusOK, nil)
}
