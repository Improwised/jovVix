package v1

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/services"
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
	presignedURLSvc       *services.PresignURLService
	logger                *zap.Logger
	event                 *events.Events
	pub                   *watermill.WatermillPublisher
}

// NewUserController returns a user
func NewUserPlayedQuizeController(goqu *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher, appConfig *config.AppConfig) (*UserPlayedQuizeController, error) {
	userPlayedQuizModel := models.InitUserPlayedQuizModel(goqu)

	activeQuizModel := models.InitActiveQuizModel(goqu, logger)

	userQuizResponseModel := models.InitUserQuizResponseModel(goqu)

	presignedURLSvc, err := services.NewFileUploadServices(appConfig.AWS.BucketName)
	if err != nil {
		return nil, err
	}

	return &UserPlayedQuizeController{
		userPlayedQuizModel:   userPlayedQuizModel,
		activeQuizModel:       activeQuizModel,
		userQuizResponseModel: userQuizResponseModel,
		presignedURLSvc:       presignedURLSvc,
		logger:                logger,
		event:                 event,
		pub:                   pub,
	}, nil
}

// ListUserPlayedQuizes to List played quiz by user
// swagger:route GET /v1/user_played_quizes UserPlayedQuiz UserPlayedQuiz
//
// List played quiz by user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseListUserPlayedQuizes
//	     401: GenericResFailConflict
//		  500: GenericResError
func (ctrl *UserPlayedQuizeController) ListUserPlayedQuizes(c *fiber.Ctx) error {
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	ctrl.logger.Debug("UserPlayedQuizeController.ListUserPlayedQuizes called", zap.Any("userId", userId))

	page := c.Query(constants.PageNumberQueryParam, "1")
	titleSearch := c.Query(constants.ParamTitle)

	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	userPlayedQuizes, count, err := ctrl.userPlayedQuizModel.ListUserPlayedQuizes(userId, pageNumber, titleSearch)
	if err != nil {
		ctrl.logger.Error("Error fetching user played quizzes", zap.Error(err))
		return err
	}

	ctrl.logger.Debug("UserPlayedQuizeController.ListUserPlayedQuizes success", zap.Any("userPlayedQuizes", userPlayedQuizes))
	return utils.JSONSuccess(c, http.StatusOK, structs.ResUserPlayedQuizWithCount{Data: userPlayedQuizes, Count: count})
}

// ListUserPlayedQuizesWithQuestionById to Analysis of played quiz by user
// swagger:route GET /v1/user_played_quizes/{user_played_quiz_id} UserPlayedQuiz RequestListUserPlayedQuizesWithQuestionById
//
// Analysis of played quiz by user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseListUserPlayedQuizesWithQuestionById
//	     401: GenericResFailConflict
//		  500: GenericResError
func (ctrl *UserPlayedQuizeController) ListUserPlayedQuizesWithQuestionById(c *fiber.Ctx) error {
	userPlayedQuizId := c.Params(constants.UserPlayedQuizId)
	ctrl.logger.Debug("UserPlayedQuizeController.ListUserPlayedQuizesWithQuestionById called", zap.Any("userPlayedQuizId", userPlayedQuizId))

	userPlayedQuizesWithQuestion, err := ctrl.userPlayedQuizModel.ListUserPlayedQuizesWithQuestionById(userPlayedQuizId)
	if err != nil {
		return err
	}

	services.ProcessAnalyticsData(userPlayedQuizesWithQuestion, ctrl.presignedURLSvc, ctrl.logger)

	ctrl.logger.Debug("UserPlayedQuizeController.ListUserPlayedQuizesWithQuestionById success", zap.Any("userPlayedQuizesWithQuestion", userPlayedQuizesWithQuestion))
	return utils.JSONSuccess(c, http.StatusOK, userPlayedQuizesWithQuestion)
}

// PlayedQuizValidation to validate the code and copy the questions in reponses for user
// swagger:route POST /v1/user_played_quizes/{invitationCode} UserPlayedQuiz RequestPlayedQuizValidation
//
//	 Validate the code and copy the questions in reponses for user.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponsePlayedQuizValidation
//		     401: GenericResFailConflict
//			  500: GenericResError
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

	response := map[string]string{
		"user_played_quiz": userPlayedQuizId.String(),
		"session_id":       session.ID.String(),
	}

	return utils.JSONSuccess(c, http.StatusOK, response)
}
