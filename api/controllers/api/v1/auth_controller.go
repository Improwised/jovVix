package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
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
	userModel *models.UserModel
	userSvc   *services.UserService
	logger    *zap.Logger
	config    config.AppConfig
}

func NewAuthController(goqu *goqu.Database, logger *zap.Logger, config config.AppConfig) (*AuthController, error) {
	userModel, err := models.InitUserModel(goqu, logger)
	if err != nil {
		return nil, err
	}

	userSvc, err := services.NewUserService(goqu, logger, config)
	if err != nil {
		return nil, err
	}

	return &AuthController{
		userModel: &userModel,
		userSvc:   userSvc,
		logger:    logger,
		config:    config,
	}, nil
}

// DoKratosAuth authenticate user with kratos session id
// swagger:route GET /v1/kratos/auth Auth DoKratosAuth
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

	return c.Redirect(ctrl.config.WebUrl)
}

// Get Details of Register user
// swagger:route GET /v1/kratos/whoami User GetRegisteredUser
//
// Get Details of Register user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseGetRegisteredUser
//	     401: GenericResFailConflict
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
	}

	ctrl.logger.Debug("AuthController.GetRegisteredUser success", zap.Any("kratosUser", kratosUser))
	return utils.JSONSuccess(c, http.StatusOK, kratosUser)
}

// Update user Details
// swagger:route PUT /v1/kratos/user User RequestUpadateRegisteredUser
//
// Update user Details.
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
//	     401: GenericResFailConflict
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

// Delete user Details
// swagger:route DELETE /v1/kratos/user User DeleteRegisteredUser
//
// Delete user Details and all its related data.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseOkWithMessage
//	     401: GenericResFailConflict
//		  500: GenericResError
func (ctrl *AuthController) DeleteRegisteredUser(c *fiber.Ctx) error {
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	kratosId := quizUtilsHelper.GetString(c.Locals(constants.KratosID))
	ctrl.logger.Debug("AuthController.DeleteRegisteredUser called", zap.Any("userID", userId), zap.Any("kratosID", kratosId))

	if kratosId == "<nil>" || userId == "<nil>" {
		ctrl.logger.Error(constants.ErrUnauthenticated)
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrUnauthenticated)
	}

	err := ctrl.userSvc.DeleteUserDataById(userId, kratosId)
	if err != nil {
		ctrl.logger.Error(constants.ErrDeleteUser, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrDeleteUser)
	}
	ctrl.logger.Debug("AuthController.DeleteRegisteredUser success", zap.Any("userID", userId), zap.Any("kratosID", kratosId))

	return utils.JSONSuccess(c, http.StatusOK, "user deleted succesfully!")
}
