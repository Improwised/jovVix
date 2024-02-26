package v1

import (
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	quiz_helper "github.com/Improwised/quizz-app/api/helpers/quiz"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type QuizController struct {
	helper *quiz_helper.HelperStructs
	logger *zap.Logger
	event  *events.Events
}

func InitQuizController(
	logger *zap.Logger,
	event *events.Events,
	pub *watermill.WatermillPublisher,
	helper *quiz_helper.HelperStructs,
) *QuizController {
	return &QuizController{
		helper: helper,
		logger: logger,
		event:  event,
	}
}

func (ctrl *QuizController) GetMyQuiz(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextUid).(string)

	quizzes, err := ctrl.helper.QuizModel.GetAllQuizzesActivity(userID)

	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)

}

func (ctrl *QuizController) CreateQuizSession(c *fiber.Ctx) error {
	userID := c.Locals(constants.ContextUid).(string)

	quizzes, err := ctrl.helper.QuizModel.GetAllQuizzesActivity(userID)

	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONSuccess(c, http.StatusOK, quizzes)

}
