package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	jwt "github.com/Improwised/quizz-app/api/pkg/jwt"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	resty "github.com/go-resty/resty/v2"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type AuthController struct {
	userService *services.UserService
	userModel   *models.UserModel
	logger      *zap.Logger
	config      config.AppConfig
}

func NewAuthController(goqu *goqu.Database, logger *zap.Logger, config config.AppConfig) (*AuthController, error) {
	userModel, err := models.InitUserModel(goqu, logger)
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(&userModel)

	return &AuthController{
		userService: userSvc,
		userModel:   &userModel,
		logger:      logger,
		config:      config,
	}, nil
}

// DoAuth authenticate user with email and password
// swagger:route POST /v1/login Auth RequestAuthnUser
//
// Authenticate user with email and password.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseAuthnUser
//		   400: GenericResFailBadRequest
//	    401: ResForbiddenRequest
//			  500: GenericResError
func (ctrl *AuthController) DoAuth(c *fiber.Ctx) error {
	var reqLoginUser structs.ReqLoginUser

	err := json.Unmarshal(c.Body(), &reqLoginUser)
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(reqLoginUser)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	user, err := ctrl.userService.Authenticate(reqLoginUser.Email, reqLoginUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.InvalidCredentials)
		}
		ctrl.logger.Error("error while get user by email and password", zap.Error(err), zap.Any("email", reqLoginUser.Email))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrLoginUser)
	}

	// token is valid for 1 hour
	token, err := jwt.CreateToken(ctrl.config, user.ID, time.Now().Add(time.Hour*1))
	if err != nil {
		ctrl.logger.Error("error while creating token", zap.Error(err), zap.Any("id", user.ID))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrLoginUser)
	}

	userCookie := &fiber.Cookie{
		Name:    constants.CookieUser,
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	}
	c.Cookie(userCookie)

	return utils.JSONSuccess(c, http.StatusOK, user)
}

// DoKratosAuth authenticate user with kratos session id
// swagger:route GET /v1/kratos/auth Auth none
//
// Authenticate user with kratos session id.
//
//			Consumes:
//			- application/json
//
//			Schemes: http, https
//		Responses:
//	      400: GenericResFailBadRequest
//		  500: GenericResError

func (ctrl *AuthController) DoKratosAuth(c *fiber.Ctx) error {
	kratosId := c.Cookies("ory_kratos_session")
	if kratosId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrKratosIDEmpty)
	}

	kratosClient := resty.New().SetBaseURL(ctrl.config.Kratos.BaseUrl+"/sessions").SetHeader("Cookie", fmt.Sprintf("%v=%v", constants.KratosCookie, kratosId)).SetHeader("accept", "application/json").SetHeader("withCredentials", "true").SetHeader("credentials", "include")

	kratosUser := config.KratosUserDetails{}
	res, err := kratosClient.R().SetResult(&kratosUser).Get("/whoami")
	if err != nil || res.StatusCode() != http.StatusOK {
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosAuth)
	}

	userStruct := models.User{}
	userStruct.KratosID = sql.NullString{String: kratosUser.Identity.ID, Valid: true}
	userStruct.CreatedAt = kratosUser.Identity.CreatedAt
	userStruct.UpdatedAt = kratosUser.Identity.UpdatedAt
	userStruct.FirstName = kratosUser.Identity.Traits.Name.First
	userStruct.LastName = kratosUser.Identity.Traits.Name.Last
	userStruct.Email = kratosUser.Identity.Traits.Email
	userStruct.Username = kratosUser.Identity.Traits.Name.First
	userStruct.Roles = "user"

	err = ctrl.userModel.InsertKratosUser(userStruct)
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
				return utils.JSONError(c, http.StatusInternalServerError, fmt.Sprintf("username (%s) already registered", userStruct.Username))
			}

			for {
				if retrying > 30 {
					break
				} else {
					userStruct.Username = quizUtilsHelper.GenerateNewStringHavingSuffixName(userStruct.Username, 5, 12)
					err = ctrl.userModel.InsertKratosUser(userStruct)
					if err != nil {
						retrying++
					} else {
						break
					}
				}
			}
		}

		if err != nil {
			ctrl.logger.Error("unable to insert the user registered with kratos into the database and username is - "+userStruct.Username, zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosDataInsertion)
		}
	}

	cookieExpirationTime, err := time.ParseDuration(ctrl.config.Kratos.CookieExpirationTime)
	if err != nil {
		ctrl.logger.Debug("unable to parse the duration for the cookie expiration", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosCookieTime)
	}

	userCookie := &fiber.Cookie{
		Name:    constants.KratosCookie,
		Value:   kratosId,
		Expires: time.Now().Add(cookieExpirationTime),
	}

	c.Cookie(userCookie)

	return c.Redirect(ctrl.config.WebUrl)
}

func (ctrl *AuthController) CreateQuickUser(c *fiber.Ctx) error {
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

func (ctrl *AuthController) IsRegisteredUser(c *fiber.Ctx) error {
	kratosId := c.Cookies("ory_kratos_session")
	if kratosId == "" {
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrKratosIDEmpty)
	}

	kratosUser := config.KratosUserDetails{}
	kratosClient := resty.New().SetBaseURL(ctrl.config.Kratos.BaseUrl+"/sessions").SetHeader("Cookie", fmt.Sprintf("%v=%v", constants.KratosCookie, kratosId)).SetHeader("accept", "application/json").SetHeader("withCredentials", "true").SetHeader("credentials", "include")

	res, err := kratosClient.R().SetResult(&kratosUser).Get("/whoami")
	if err != nil || res.StatusCode() != http.StatusOK {
		ctrl.logger.Debug("unauthenticated registration", zap.Any("response from kratos", res.RawResponse), zap.Error(err), zap.Any("kratos response", res))
		return utils.JSONError(c, res.StatusCode(), constants.ErrKratosAuth)
	} else {
		return utils.JSONSuccess(c, http.StatusOK, true)
	}
}

func (ctrl *AuthController) GetRegisteredUser(c *fiber.Ctx) error {
	kratosId := c.Cookies("ory_kratos_session")
	ctrl.logger.Debug("AuthController.GetRegisteredUser called", zap.Any("ory_kratos_session", kratosId))
	if kratosId == "" {
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrKratosIDEmpty)
	}

	kratosUser := config.KratosUserDetails{}
	kratosClient := resty.New().SetBaseURL(ctrl.config.Kratos.BaseUrl+"/sessions").SetHeader("Cookie", fmt.Sprintf("%v=%v", constants.KratosCookie, kratosId)).SetHeader("accept", "application/json").SetHeader("withCredentials", "true").SetHeader("credentials", "include")

	res, err := kratosClient.R().SetResult(&kratosUser).Get("/whoami")
	if err != nil || res.StatusCode() != http.StatusOK {
		ctrl.logger.Debug("unauthenticated registration", zap.Any("response from kratos", res.RawResponse), zap.Error(err), zap.Any("kratos response", res))
		return utils.JSONError(c, res.StatusCode(), constants.ErrKratosAuth)
	} else {
		ctrl.logger.Debug("AuthController.GetRegisteredUser success", zap.Any("kratosUser", kratosUser))
		return utils.JSONSuccess(c, http.StatusOK, kratosUser)
	}
}

func (ctrl *AuthController) UpadateRegisteredUser(c *fiber.Ctx) error {
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	ctrl.logger.Debug("AuthController.UpadateRegisteredUser called", zap.Any("userId", userId))

	ctrl.logger.Debug("validate req", zap.Any("Body", c.Body()))
	var userReq structs.ReqUpdateUser
	err := json.Unmarshal(c.Body(), &userReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(userReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}
	ctrl.logger.Debug("validate req success", zap.Any("userReq", userReq))

	ctrl.logger.Debug("userModel.IsUniqueEmail called", zap.Any("Email", userReq.Email))
	IsUniqueEmail, err := ctrl.userModel.IsUniqueEmailExceptId(userId, userReq.Email)
	if err != nil {
		ctrl.logger.Error("error while UpdateUser", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, err.Error())
	}
	ctrl.logger.Debug("userModel.IsUniqueEmail success", zap.Any("IsUniqueEmail", IsUniqueEmail))

	if !IsUniqueEmail {
		return utils.JSONFail(c, http.StatusBadRequest, "email already exist")
	}

	var kratosStruct = config.KratosTraits{}
	kratosStruct.Name.First = userReq.FirstName
	kratosStruct.Name.Last = userReq.LastName
	kratosStruct.Email = userReq.Email

	kratosjson, err := json.Marshal(kratosStruct)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	ctrl.logger.Debug("userModel.UpdateKratosUserDetails called", zap.Any("userId", userId))
	if err := ctrl.userModel.UpdateKratosUserDetails(models.User{ID: userId, FirstName: userReq.FirstName, LastName: userReq.LastName, Email: userReq.Email}, kratosjson); err != nil {
		ctrl.logger.Error("error while UpdateKratosUserById", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, err.Error())
	}
	ctrl.logger.Debug("userModel.UpdateKratosUserDetails success", zap.Any("userId", userId))

	ctrl.logger.Debug("userModel.GetById success", zap.Any("userId", userId))
	user, err := ctrl.userModel.GetById(userId)
	if err != nil {
		ctrl.logger.Error("error while GetById", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, err.Error())
	}
	ctrl.logger.Debug("userModel.GetById success", zap.Any("user", user))

	return utils.JSONSuccess(c, http.StatusOK, user)
}
