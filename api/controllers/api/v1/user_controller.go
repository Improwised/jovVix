package v1

import (
	"database/sql"
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserController for user controllers
type UserController struct {
	userService *services.UserService
	userModel   *models.UserModel
	logger      *zap.Logger
	event       *events.Events
	pub         *watermill.WatermillPublisher
}

// NewUserController returns a user
func NewUserController(goqu *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher) (*UserController, error) {
	userModel, err := models.InitUserModel(goqu, logger)
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(&userModel)

	return &UserController{
		userService: userSvc,
		userModel:   &userModel,
		logger:      logger,
		event:       event,
		pub:         pub,
	}, nil
}

// GetUserMeta Get Details of user
// swagger:route GET /v1/user/who User GetUserMeta
//
// Get Details of user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseUserDetails
//	     401: GenericResFailConflict
//		  500: GenericResError
func (ctrl *UserController) GetUserMeta(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	kratosID := quizUtilsHelper.GetString(c.Locals(constants.KratosID))
	ctrl.logger.Debug("UserController.GetUserMeta called", zap.Any("userID", userID), zap.Any("kratosID", kratosID))

	if kratosID == "<nil>" && userID == "<nil>" {
		ctrl.logger.Error(constants.ErrUnauthenticated)
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrUnauthenticated)
	}

	if kratosID != "<nil>" {
		user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
		if !ok {
			ctrl.logger.Error("Cannot be able to get the userMeta details from database")
			return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrGetUser)
		}

		ctrl.logger.Debug("UserController.GetUserMeta success", zap.Any("user", user))
		return utils.JSONSuccess(c, http.StatusOK, map[string]string{
			"username":  user.Username,
			"firstname": user.FirstName,
			"email":     user.Email,
			"role":      "admin-user",
			"avatar":    user.ImageKey,
		})
	}

	ctrl.logger.Debug("userModel.GetById called", zap.Any("userID", userID))
	user, err := ctrl.userModel.GetById(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctrl.logger.Error(constants.ErrGetUser, zap.Error(err))
			return utils.JSONError(c, http.StatusNotFound, constants.ErrGetUser)
		}
		ctrl.logger.Error(constants.ErrGetUser, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
	}
	ctrl.logger.Debug("userModel.GetById success", zap.Any("user", user))
	ctrl.logger.Debug("UserController.GetUserMeta success", zap.Any("user", user))

	return utils.JSONSuccess(c, http.StatusOK, map[string]string{
		"username":  user.Username,
		"firstname": user.FirstName,
		"email":     user.Email,
		"role":      "guest-user",
		"avatar":    user.ImageKey,
	})
}
