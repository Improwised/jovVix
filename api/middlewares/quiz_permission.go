package middlewares

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// check for read permission of perticular quiz and also save permission in context
func (m *Middleware) QuizPermission(c *fiber.Ctx) error {
	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
	if !ok {
		m.Logger.Error("error while fetching user context from connection")
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrUnauthenticated)
	}

	quizId := c.Params(constants.QuizId)
	if quizId == "" {
		return utils.JSONError(c, http.StatusBadRequest, "No quiz_id found")
	}

	isQuizCreator, err := m.sharedQuizzesModel.CheckQuizCreatorExists(quizId, user.ID)
	if err != nil {
		m.Logger.Error(constants.ErrGetUser, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrCheckQuizCreatorExists)
	}

	if isQuizCreator {
		c.Locals(constants.ContextQuizPermission, constants.SharePermission)
		return c.Next()
	}

	permission, err := m.sharedQuizzesModel.GetPermissionByQuizAndUser(quizId, user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			m.Logger.Error(constants.ErrNotAllowed, zap.Error(err))
			return utils.JSONError(c, http.StatusUnauthorized, constants.ErrNotAllowed)
		}
		m.Logger.Error(constants.ErrGetQuizPermission, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetQuizPermission)
	}

	c.Locals(constants.ContextQuizPermission, permission)
	return c.Next()
}

// check for edit permission of perticular quiz
func (m *Middleware) VerifyQuizEditAccess(c *fiber.Ctx) error {
	// Retrieve the user's permission for the current quiz from context
	permission := fmt.Sprintf("%s", c.Locals(constants.ContextQuizPermission))
	if permission != constants.SharePermission && permission != constants.WritePermission {
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
	}

	return c.Next()
}

// check for share permission of perticular quiz
func (m *Middleware) VerifyQuizShareAccess(c *fiber.Ctx) error {
	// Retrieve the user's permission for the current quiz from context
	permission := fmt.Sprintf("%s", c.Locals(constants.ContextQuizPermission))
	if permission != constants.SharePermission {
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrUnauthorized)
	}

	return c.Next()
}
