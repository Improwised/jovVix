package v1

import (
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserController for user controllers
type UserPlayedQuizeController struct {
	userPlayedQuizModel *models.UserPlayedQuizModel
	logger              *zap.Logger
	event               *events.Events
	pub                 *watermill.WatermillPublisher
}

// NewUserController returns a user
func NewUserPlayedQuizeController(goqu *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher) (*UserPlayedQuizeController, error) {
	userPlayedQuizModel := models.InitUserPlayedQuizModel(goqu)

	return &UserPlayedQuizeController{
		userPlayedQuizModel: userPlayedQuizModel,
		logger:              logger,
		event:               event,
		pub:                 pub,
	}, nil
}

func (ctrl *UserPlayedQuizeController) ListUserPlayedQuizes(c *fiber.Ctx) error {
	userId := c.Locals(constants.ContextUid).(string)
	ctrl.logger.Debug("UserPlayedQuizeController.ListUserPlayedQuizes called", zap.Any("userId", userId))

	userPlayedQuizes, err := ctrl.userPlayedQuizModel.ListUserPlayedQuizes(userId)
	if err != nil {
		return err
	}

	ctrl.logger.Debug("UserPlayedQuizeController.ListUserPlayedQuizes success", zap.Any("userPlayedQuizes", userPlayedQuizes))
	return utils.JSONSuccess(c, http.StatusOK, userPlayedQuizes)
}

func (ctrl *UserPlayedQuizeController) ListUserPlayedQuizesWithQuestionById(c *fiber.Ctx) error {
	userPlayedQuizId := c.Params(constants.UserPlayedQuizId)
	ctrl.logger.Debug("UserPlayedQuizeController.ListUserPlayedQuizesWithQuestionById called", zap.Any("userPlayedQuizId", userPlayedQuizId))

	userPlayedQuizesWithQuestion, err := ctrl.userPlayedQuizModel.ListUserPlayedQuizesWithQuestionById(userPlayedQuizId)
	if err != nil {
		return err
	}

	ctrl.logger.Debug("UserPlayedQuizeController.ListUserPlayedQuizesWithQuestionById success", zap.Any("userPlayedQuizesWithQuestion", userPlayedQuizesWithQuestion))
	return utils.JSONSuccess(c, http.StatusOK, userPlayedQuizesWithQuestion)
}
