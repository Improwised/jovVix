package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/pkg/templates"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type SharedQuizzes struct {
	sharedQuizzesModel *models.SharedQuizzesModel
	userModel          *models.UserModel
	emailService       *services.EmailService
	logger             *zap.Logger
	config             *config.AppConfig
}

func NewSharedQuizzesController(goqu *goqu.Database, logger *zap.Logger, config *config.AppConfig) (*SharedQuizzes, error) {

	sharedQuizzesModel := models.InitSharedQuizzesModel(goqu, logger)
	userModel, err := models.InitUserModel(goqu, logger)
	if err != nil {
		return nil, err
	}

	emailService := services.NewEmailService(logger, &config.SMTP)

	return &SharedQuizzes{
		sharedQuizzesModel: sharedQuizzesModel,
		userModel:          &userModel,
		emailService:       emailService,
		logger:             logger,
		config:             config,
	}, nil
}

// ShareQuiz to insert data for share the quiz.
// swagger:route POST /v1/shared_quizzes ShareQuiz RequestShareQuiz
//
// share quiz to other user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseOkWithMessage
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (sqctrl *SharedQuizzes) ShareQuiz(c *fiber.Ctx) error {
	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
	if !ok {
		sqctrl.logger.Error(constants.ErrConvertTypeUser)
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrShareQuiz)
	}

	quizId := c.Params(constants.QuizId)
	sqctrl.logger.Debug("SharedQuizzes.ShareQuiz called", zap.Any("quizId", quizId), zap.Any("userId", user.ID))
	if quizId == "" {
		return utils.JSONError(c, http.StatusBadRequest, "No quiz_id found")
	}

	sqctrl.logger.Debug("validate req", zap.Any("Body", c.Body()))
	var shareQuizReq structs.ReqShareQuiz
	err := json.Unmarshal(c.Body(), &shareQuizReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(shareQuizReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}
	sqctrl.logger.Debug("validate req success", zap.Any("shareQuizReq", shareQuizReq))

	// Insert for share quiz
	id, err := sqctrl.sharedQuizzesModel.InsertSharedQuiz(models.SharedQuizzes{
		QuizId:     quizId,
		SharedTo:   shareQuizReq.Email,
		SharedBy:   user.ID,
		Permission: shareQuizReq.Permission,
	})
	if err != nil {
		sqctrl.logger.Error(constants.ErrShareQuiz, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrShareQuiz)
	}

	// Send Email to user to notify
	quizLink := fmt.Sprintf("%s/admin/quiz/list-quiz/%s", sqctrl.config.WebUrl, quizId)
	emailBody := templates.GenerateQuizShareEmail(quizLink, user.Email, shareQuizReq.Permission)

	err = sqctrl.emailService.SendEmail(shareQuizReq.Email, constants.QuizEmailSubject, emailBody)
	if err != nil {
		sqctrl.logger.Error("Failed to send email", zap.Error(err))
	}
	sqctrl.logger.Debug("SharedQuizzes.ShareQuiz success", zap.Any("quizId", quizId), zap.Any("userId", user.ID))

	return utils.JSONSuccess(c, http.StatusOK, id)
}

// ListQuizAuthorizedUsers to List authorized users for perticular quiz.
// swagger:route GET /v1/shared_quizzes/{quiz_id} ShareQuiz RequestListQuizAuthorizedUsers
//
// List authorized users for perticular quiz.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseListQuizAuthorizedUsers
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (sqctrl *SharedQuizzes) ListQuizAuthorizedUsers(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)
	sqctrl.logger.Debug("SharedQuizzes.ListQuizAuthorizedUsers called", zap.Any("quizId", quizId))
	if quizId == "" {
		return utils.JSONError(c, http.StatusBadRequest, "No quiz_id found")
	}

	quizAuthorizedUsers, err := sqctrl.sharedQuizzesModel.ListQuizAuthorizedUsersByQuizId(quizId)
	if err != nil {
		sqctrl.logger.Error(constants.ErrFetchAuthorizedUsersError, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrFetchAuthorizedUsersError)
	}
	sqctrl.logger.Debug("SharedQuizzes.ListQuizAuthorizedUsers success", zap.Any("quizAuthorizedUsers", quizAuthorizedUsers))

	return utils.JSONSuccess(c, http.StatusOK, quizAuthorizedUsers)
}

// UpdateUserPermissionOfQuiz to Update authorized user permission for perticular quiz.
// swagger:route PUT /v1/shared_quizzes/{quiz_id} ShareQuiz RequestUpdateUserPermissionOfQuiz
//
// Update authorized user permission for perticular quiz.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseOkWithMessage
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (sqctrl *SharedQuizzes) UpdateUserPermissionOfQuiz(c *fiber.Ctx) error {
	sharedQuizId := c.Query(constants.SharedQuizId)
	sqctrl.logger.Debug("SharedQuizzes.UpdateUserPermissionOfQuiz called", zap.Any(constants.SharedQuizId, sharedQuizId))
	if sharedQuizId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.BadRequestSharedQuizIdNotFound)
	}

	sqctrl.logger.Debug("validate req", zap.Any("Body", c.Body()))
	var shareQuizReq structs.ReqShareQuiz
	err := json.Unmarshal(c.Body(), &shareQuizReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(shareQuizReq)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}
	sqctrl.logger.Debug("validate req success", zap.Any("shareQuizReq", shareQuizReq))

	err = sqctrl.sharedQuizzesModel.UpdateUserPermissionById(sharedQuizId, shareQuizReq)
	if err != nil {
		sqctrl.logger.Error(constants.ErrUpdateUserPermissionForQuiz, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUpdateUserPermissionForQuiz)
	}
	sqctrl.logger.Debug("SharedQuizzes.UpdateUserPermissionOfQuiz success", zap.Any("shareQuizReq", shareQuizReq))

	return utils.JSONSuccess(c, http.StatusOK, "User permission updated successfully!")
}

// DeleteUserPermissionOfQuiz to Delete authorized user permission for perticular quiz.
// swagger:route DELETE /v1/shared_quizzes/{quiz_id} ShareQuiz RequestDeleteUserPermissionOfQuiz
//
// Delete authorized user permission for perticular quiz.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseOkWithMessage
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (sqctrl *SharedQuizzes) DeleteUserPermissionOfQuiz(c *fiber.Ctx) error {
	sharedQuizId := c.Query(constants.SharedQuizId)
	sqctrl.logger.Debug("SharedQuizzes.DeleteUserPermissionOfQuiz called", zap.Any("sharedQuizId", sharedQuizId))
	if sharedQuizId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.BadRequestSharedQuizIdNotFound)
	}

	err := sqctrl.sharedQuizzesModel.DeleteUserPermissionById(sharedQuizId)
	if err != nil {
		sqctrl.logger.Error(constants.ErrDeleteUserPermissionForQuiz, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrDeleteUserPermissionForQuiz)
	}
	sqctrl.logger.Debug("SharedQuizzes.DeleteUserPermissionOfQuiz success", zap.Any("sharedQuizId", sharedQuizId))

	return utils.JSONSuccess(c, http.StatusOK, "User permission deleted successfully!")
}

// ListSharedQuizzes to List shared quiz for perticular user (only shared with the user or shared by the user).
// swagger:route GET /v1/shared_quizzes ShareQuiz RequestListSharedQuizzes
//
// List shared quiz for perticular user (only shared with the user or shared by the user).
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseListSharedQuizzes
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (sqctrl *SharedQuizzes) ListSharedQuizzes(c *fiber.Ctx) error {
	requestType := c.Query("type")
	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
	if !ok {
		sqctrl.logger.Error(constants.ErrConvertTypeUser)
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrListShareQuiz)
	}
	sqctrl.logger.Debug("SharedQuizzes.ListSharedQuizzes called", zap.Any("userId", user.ID), zap.Any("requestType", requestType))

	var sharedTo, sharedBy string
	switch requestType {
	case "shared_by_me":
		// Quizzes shared by the user
		sharedBy = user.ID
	case "shared_with_me":
		// Quizzes shared with the user
		sharedTo = user.Email
	default:
		return utils.JSONError(c, http.StatusBadRequest, "Invalid request type. Use 'shared_by_me' or 'shared_with_me'.")
	}

	sharedQuizzes, err := sqctrl.sharedQuizzesModel.ListSharedQuizzes(sharedBy, sharedTo)
	if err != nil {
		sqctrl.logger.Error(constants.ErrListShareQuiz, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrListShareQuiz)
	}
	sqctrl.logger.Debug("SharedQuizzes.ListSharedQuizzes success", zap.Any("sharedQuizzes", sharedQuizzes))

	return utils.JSONSuccess(c, http.StatusOK, sharedQuizzes)
}
