package middlewares

import (
	"database/sql"
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/utils"
	fiber "github.com/gofiber/fiber/v2"
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

func (rpm *RolePermissionMiddleware) IsAllowed(c *fiber.Ctx) error {

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
			return utils.JSONFail(c, http.StatusNotFound, constants.ErrUnauthenticated)
		}

		userID := quizUtilsHelper.GetString(userAny)
		user, err := rpm.middleware.UserService.GetUser(userID)

		if err != nil {
			if err == sql.ErrNoRows {
				return utils.JSONFail(c, http.StatusNotFound, constants.ErrUnauthenticated)
			}

			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
		}

		c.Locals(constants.ContextUser, user)
		userLocal = user
	}

	if !rpm.allowedRoles.IsAllowed(models.Role((userLocal.Roles))) {
		return utils.JSONFail(c, http.StatusUnauthorized, constants.ErrNotAllowed)
	}

	return c.Next()
}
