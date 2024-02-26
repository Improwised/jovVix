package middlewares

import (
	"database/sql"
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
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
	userLocal := c.Locals(constants.ContextUser)

	if userLocal == any(nil) {
		userID := c.Locals(constants.ContextUid)
		user, err := rpm.middleware.UserService.GetUser(userID.(string))

		if err != nil {
			if err == sql.ErrNoRows {
				return utils.JSONError(c, http.StatusNotFound, constants.ErrUnauthenticated)
			}

			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
		}

		c.Locals(constants.ContextUser, user)
		userLocal = user
	}

	if !rpm.allowedRoles.IsAllowed(models.Role(userLocal.(models.User).Roles)) {
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrNotAllowed)
	}

	return c.Next()
}
