package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Improwised/quizz-app/api/cli/workers"
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
	"github.com/lib/pq"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

// UserController for user controllers
type UserController struct {
	userService *services.UserService
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
		logger:      logger,
		event:       event,
		pub:         pub,
	}, nil
}

// UserGet get the user by id
// swagger:route GET /v1/users/{userId} Users RequestGetUser
//
// Get a user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseGetUser
//	   400: GenericResFailNotFound
//		  500: GenericResError
func (ctrl *UserController) GetUser(c *fiber.Ctx) error {
	userID := c.Params(constants.ParamUid)
	user, err := ctrl.userService.GetUser(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctrl.logger.Debug("No user found", zap.Error(err))
			return utils.JSONFail(c, http.StatusNotFound, constants.UserNotExist)
		}
		ctrl.logger.Error("error while get user by id", zap.Any("id", userID), zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
	}
	return utils.JSONSuccess(c, http.StatusOK, user)
}

// CreateUser registers a user
// swagger:route POST /v1/users Users RequestCreateUser
//
// Register a user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  201: ResponseCreateUser
//	   400: GenericResFailBadRequest
//		  500: GenericResError
func (ctrl *UserController) CreateUser(c *fiber.Ctx) error {

	var userReq structs.ReqRegisterUser

	err := json.Unmarshal(c.Body(), &userReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(userReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	if len(userReq.UserName) > 12 || len(userReq.UserName) < 6 {
		return utils.JSONFail(c, http.StatusBadRequest, "username should be in between 6 to 12 character length")
	}

	role := "user"

	if userReq.Password == "" {
		return utils.JSONFail(c, http.StatusBadRequest, "password must not be empty")
	}
	user, err := ctrl.userService.RegisterUser(models.User{
		FirstName: userReq.FirstName,
		LastName:  userReq.LastName,
		Email:     userReq.Email,
		Password: sql.NullString{
			String: userReq.Password,
			Valid:  true,
		},
		Roles:    role,
		Username: userReq.UserName}, ctrl.event)
	if err != nil {

		if err.(*pq.Error).Constraint == constants.UserUkey {
			return utils.JSONError(c, http.StatusBadRequest, constants.ErrUsernameExists)
		}

		ctrl.logger.Error("error while insert user", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrInsertUser)
	}

	// publish message to queue
	welcomeMail := workers.WelcomeMail{FirstName: userReq.FirstName, LastName: userReq.LastName, Email: userReq.Email, Roles: role}
	err = ctrl.pub.Publish("user", welcomeMail)
	if err != nil {
		ctrl.logger.Error("error while publish message", zap.Error(err))
	}

	return utils.JSONSuccess(c, http.StatusCreated, user)
}

func (ctrl *UserController) IsAdmin(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	user, err := ctrl.userService.GetUser(userID)

	aqi := c.Query("active_quiz_id")

	if err != nil {
		if err == sql.ErrNoRows {
			RemoveUserToken(constants.ContextUid)
			return utils.JSONFail(c, http.StatusNotFound, constants.UserNotExist)
		}
		return utils.JSONError(c, http.StatusBadRequest, constants.UnknownError)
	}

	if user.Roles == "admin" {
		if aqi != "" {
			return c.Next()
		}
		return utils.JSONSuccess(c, http.StatusOK, true)
	}

	return utils.JSONFail(c, http.StatusBadRequest, constants.Unauthenticated)
}

func (ctrl *UserController) GetUserMeta(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	var userMeta = models.User{}
	var err error
	var ok bool

	if kratosToken := c.Cookies(constants.KratosCookie); kratosToken == "" {
		userMeta, err = ctrl.userService.GetUser(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Cookie(RemoveUserToken(constants.ContextUid))
				ctrl.logger.Error("Cannot be able to get the userMeta details from database")
				return utils.JSONFail(c, http.StatusNotFound, constants.Unauthenticated)
			}
			return utils.JSONError(c, http.StatusBadRequest, constants.UnknownError)
		}
	} else {
		userMeta, ok = quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))

		if !ok {
			ctrl.logger.Error("Cannot be able to get the userMeta details from database")
			return utils.JSONFail(c, http.StatusNotFound, constants.Unauthenticated)
		}
	}

	return utils.JSONSuccess(c, http.StatusOK, map[string]string{
		"username":  userMeta.Username,
		"firstname": userMeta.FirstName,
		"email":     userMeta.Email,
	})
}

func RemoveUserToken(key string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.Value = "" // Or set a generic value like "deleted"
	cookie.Expires = time.Now().Add(-1 * time.Second)
	return cookie
}
