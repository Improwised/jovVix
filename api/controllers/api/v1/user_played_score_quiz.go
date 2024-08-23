package v1

import (
	"database/sql"
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
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
	userPlayedQuizModel   *models.UserPlayedQuizModel
	activeQuizModel       *models.ActiveQuizModel
	userQuizResponseModel *models.UserQuizResponseModel
	logger                *zap.Logger
	event                 *events.Events
	pub                   *watermill.WatermillPublisher
}

// NewUserController returns a user
func NewUserPlayedQuizeController(goqu *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher) (*UserPlayedQuizeController, error) {
	userPlayedQuizModel := models.InitUserPlayedQuizModel(goqu)

	activeQuizModel := models.InitActiveQuizModel(goqu, logger)

	userQuizResponseModel := models.InitUserQuizResponseModel(goqu)

	return &UserPlayedQuizeController{
		userPlayedQuizModel:   userPlayedQuizModel,
		activeQuizModel:       activeQuizModel,
		userQuizResponseModel: userQuizResponseModel,
		logger:                logger,
		event:                 event,
		pub:                   pub,
	}, nil
}

func (ctrl *UserPlayedQuizeController) ListUserPlayedQuizes(c *fiber.Ctx) error {
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
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

func (ctrl *UserPlayedQuizeController) PlayedQuizValidation(c *fiber.Ctx) error {
	invitationCode := c.Params(constants.QuizSessionInvitationCode)
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	ctrl.logger.Debug("UserPlayedQuizeController.PlayedQuizValidation called", zap.Any("userId", userId), zap.Any("invitationCode", invitationCode))

	if userId == "<nil>" {
		ctrl.logger.Error("error while PlayedQuizValidation")
		return utils.JSONFail(c, http.StatusInternalServerError, "userId is not found")
	}

	ctrl.logger.Debug("activeQuizModel.GetSessionByCode called", zap.Any("userId", userId))
	session, err := ctrl.activeQuizModel.GetSessionByCode(invitationCode)
	if err != nil {
		if err == sql.ErrNoRows {
			ctrl.logger.Error(constants.ErrInvitationCodeNotFound, zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrInvitationCodeNotFound)
		}
		ctrl.logger.Error("error in invitation code", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, "error while get session by invitation code")
	}
	ctrl.logger.Debug("activeQuizModel.activeQuizModel success", zap.Any("session", session))

	if userId == session.AdminID {
		ctrl.logger.Error(constants.ErrAdminCannotBeUser)
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrAdminCannotBeUser)
	}

	ctrl.logger.Debug("userPlayedQuizModel.CreateUserPlayedQuizIfNotExists called", zap.Any("userId", userId), zap.Any("sessionID", session.ID))
	userPlayedQuizId, isNonExistingParticipants, err := ctrl.userPlayedQuizModel.CreateUserPlayedQuizIfNotExists(userId, session.ID)
	if err != nil {
		ctrl.logger.Error(constants.ErrUserQuizSessionValidation, zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrUserQuizSessionValidation)
	}
	ctrl.logger.Debug("userPlayedQuizModel.CreateUserPlayedQuizIfNotExists success", zap.Any("userPlayedQuizId", userPlayedQuizId))

	// insert initial records of question if user is new for quiz(isNonExistingParticipants == 2 => new user)
	if isNonExistingParticipants == 2 {
		ctrl.logger.Debug("userQuizResponseModel.GetQuestionsCopy called", zap.Any("userPlayedQuizId", userPlayedQuizId), zap.Any("sessionID", session.ID))
		err = ctrl.userQuizResponseModel.GetQuestionsCopy(userPlayedQuizId, session.ID)
		if err != nil {
			ctrl.logger.Error("error while get question copy", zap.Error(err))
			return utils.JSONFail(c, http.StatusInternalServerError, "error while get question copy")
		}
		ctrl.logger.Debug("userQuizResponseModel.GetQuestionsCopy success")
	}

	ctrl.logger.Debug("UserPlayedQuizeController.PlayedQuizValidation success", zap.Any("userPlayedQuizId", userPlayedQuizId.String()))

	return utils.JSONSuccess(c, http.StatusOK, userPlayedQuizId.String())
}
