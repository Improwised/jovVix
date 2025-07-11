package v1

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/jwt"
	"github.com/Improwised/jovvix/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

// UserController for user controllers
type UserController struct {
	userModel *models.UserModel
	logger    *zap.Logger
	config    config.AppConfig
}

// NewUserController returns a user
func NewUserController(goqu *goqu.Database, logger *zap.Logger, config config.AppConfig) (*UserController, error) {
	userModel, err := models.InitUserModel(goqu, logger)
	if err != nil {
		return nil, err
	}

	return &UserController{
		userModel: &userModel,
		logger:    logger,
		config:    config,
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

// Create Guest user to play for quiz directly without login
// swagger:route POST /v1/user/{username} User RequestCreateQuickUser
//
// Create Guest user to play for quiz directly without login.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseUserDetails
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (ctrl *UserController) CreateGuestUser(c *fiber.Ctx) error {
	username := c.Params(constants.Username)

	avatarName := c.Query("avatar_name")
	if username == "" || avatarName == "" {
		return utils.JSONError(c, http.StatusBadRequest, "please provide username and avatar name")
	}

	userObj := models.User{
		FirstName: username,
		Username:  username,
		Roles:     "user",
		ImageKey:  avatarName,
	}
	user, err := ctrl.userModel.InsertUser(userObj)
	if err != nil {
		pqErr, ok := quizUtilsHelper.ConvertType[*pq.Error](err)
		retrying := 0

		if !ok {
			ctrl.logger.Debug("unable to convert postgres error")
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrorTypeConversion)
		}

		if pqErr.Code == "23505" {

			if pqErr.Constraint != constants.UserUkey {
				ctrl.logger.Debug("user wants to use the username which is already registered")
				return utils.JSONError(c, http.StatusInternalServerError, fmt.Sprintf("username (%s) already registered", user.Username))
			}

			for {
				if retrying > 30 {
					break
				} else {
					user.Username = quizUtilsHelper.GenerateNewStringHavingSuffixName(user.Username, 5, 12)
					user, err = ctrl.userModel.InsertUser(user)
					if err != nil {
						retrying++
					} else {
						break
					}
				}
			}
		}

		if err != nil {
			ctrl.logger.Error("unable to insert the user registered with kratos into the database and username is - "+user.Username, zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosDataInsertion)
		}
	}

	cookieExpirationTime, err := time.ParseDuration(ctrl.config.Kratos.CookieExpirationTime)
	if err != nil {
		ctrl.logger.Debug("unable to parse the duration for the cookie expiration", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosCookieTime)
	}

	token, err := jwt.CreateToken(ctrl.config, user.ID, time.Now().Add(time.Hour*2))
	if err != nil {
		ctrl.logger.Error("error while creating token", zap.Error(err), zap.Any("id", user.ID))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrLoginUser)
	}

	userCookie := &fiber.Cookie{
		Name:    constants.CookieUser,
		Value:   token,
		Expires: time.Now().Add(cookieExpirationTime),
	}

	c.Cookie(userCookie)

	return utils.JSONSuccess(c, http.StatusOK, user)
}
