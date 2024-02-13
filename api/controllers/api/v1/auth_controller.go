package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/models"
	jwt "github.com/Improwised/quizz-app/api/pkg/jwt"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	resty "github.com/go-resty/resty/v2"
	fiber "github.com/gofiber/fiber/v2"
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
	userModel, err := models.InitUserModel(goqu)
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
// swagger:route POST /login Auth RequestAuthnUser
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
// swagger:route GET /kratos/auth Auth none
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
	kratosID := c.Locals(constants.KratosID)

	if kratosID.(string) == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrKratosIDEmpty)
	}

	kratosClient := resty.New().SetBaseURL(ctrl.config.Kratos.BaseUrl+"/sessions").SetHeader("Cookie", fmt.Sprintf("%v=%v", constants.KratosCookie, kratosID)).SetHeader("accept", "application/json")

	kratosUser := config.KratosUserDetails{}
	res, err := kratosClient.R().SetResult(&kratosUser).Get("/whoami")
	if err != nil || res.StatusCode() != http.StatusOK {
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosAuth)
	}

	userStruct := models.User{}
	userStruct.KratosID = kratosUser.Identity.ID
	userStruct.FirstName = kratosUser.Identity.Traits.Name.First
	userStruct.LastName = kratosUser.Identity.Traits.Name.Last
	userStruct.Email = kratosUser.Identity.Traits.Email
	userStruct.CreatedAt = kratosUser.Identity.CreatedAt
	userStruct.UpdatedAt = kratosUser.Identity.UpdatedAt

	err = ctrl.userModel.InsertKratosUser(userStruct)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosDataInsertion)
	}

	cookieExpirationTime, err := time.ParseDuration(ctrl.config.Kratos.CookieExpirationTime)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosCookieTime)
	}

	userCookie := &fiber.Cookie{
		Name:    constants.KratosCookie,
		Value:   kratosID.(string),
		Expires: time.Now().Add(cookieExpirationTime),
	}

	c.Cookie(userCookie)

	return c.Redirect(ctrl.config.Kratos.UIUrl)
}
