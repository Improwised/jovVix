package middlewares

import (
	"database/sql"
	"fmt"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type RolePermissionMiddleware struct {
	middleware   Middleware
	allowedRoles models.AllowedRoles
}

func NewRolePermissionMiddleware(middleware Middleware, allowedRoles models.AllowedRoles) RolePermissionMiddleware {
	return RolePermissionMiddleware{
		middleware:   middleware,
		allowedRoles: allowedRoles,
	}
}

func (m *RolePermissionMiddleware) IsAllowed(c *fiber.Ctx) error {

	if c.Locals(constants.MiddlewareError) != nil {
		return c.Next()
	}

	userAny := c.Locals(constants.ContextUser)
	userLocal, ok := quizUtilsHelper.ConvertType[models.User](userAny)

	// if user object not exists then check for userID
	if !ok {
		userAny = c.Locals(constants.ContextUid)

		// if userID not exists then take it as fail
		if userAny == any(nil) {
			c.Locals(constants.MiddlewareError, constants.ErrUnauthenticated)
			m.middleware.Logger.Error("Username not provided", zap.Error(fmt.Errorf(constants.ErrUserRequiredToCheckRole)))
			return c.Next()
		}

		userID := quizUtilsHelper.GetString(userAny)
		user, err := m.middleware.UserService.GetUser(userID)

		if err != nil {
			if err == sql.ErrNoRows {
				c.Locals(constants.MiddlewareError, constants.ErrUnauthenticated)
				m.middleware.Logger.Error("User not found", zap.Error(fmt.Errorf(constants.UserNotExist)))
				return c.Next()
			}

			c.Locals(constants.MiddlewareError, constants.ErrGetUser)
			m.middleware.Logger.Error("Error in Get user", zap.Error(fmt.Errorf(constants.ErrGetUser)))
			return c.Next()
		}

		c.Locals(constants.ContextUser, user)
		userLocal = user
	}

	if !m.allowedRoles.IsAllowed(models.Role((userLocal.Roles))) {
		c.Locals(constants.MiddlewareError, constants.ErrNotAllowed)
		m.middleware.Logger.Error("User have no demanded role", zap.Error(fmt.Errorf(constants.ErrNotAllowed)))
	}

	return c.Next()
}
