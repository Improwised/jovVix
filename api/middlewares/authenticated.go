package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/pkg/jwt"
	"github.com/Improwised/quizz-app/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	j "github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

func (m *Middleware) Authenticated(c *fiber.Ctx) error {

	if m.Config.Kratos.IsEnabled {
		sessionID := c.Cookies("ory_kratos_session")
		if sessionID == "" {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}
		c.Locals(constants.KratosID, sessionID)
		return c.Next()
	}

	token := c.Cookies(constants.CookieUser, "")
	if token == "" {
		fmt.Println("no token***********************")
		return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
	}

	claims, err := jwt.ParseToken(m.Config, token)
	if err != nil {
		fmt.Println("no parser token********************")
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}

		m.Logger.Error("error while checking user identity", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthenticated)
	}

	c.Locals(constants.ContextUid, claims.Subject())
	return c.Next()
}
