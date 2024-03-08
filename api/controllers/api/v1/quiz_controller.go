package v1

import (
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	quiz_helper "github.com/Improwised/quizz-app/api/helpers/quiz"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

	if currentQuiz == "" {
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuizNotFound)
	}

	return utils.JSONSuccess(c, http.StatusAccepted, nil)
}
