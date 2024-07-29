package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/jwt"
	"github.com/Improwised/quizz-app/api/utils"
	resty "github.com/go-resty/resty/v2"
	fiber "github.com/gofiber/fiber/v2"
	j "github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

func (m *Middleware) Authenticated(c *fiber.Ctx) error {
	token := c.Cookies(constants.CookieUser, "")
	if token == "" {
		m.Logger.Debug("returning unauthorized after found empty token")
		return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
	}

	claims, err := jwt.ParseToken(m.Config, token)
	if err != nil {
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			c.Cookie(RemoveCookie(constants.CookieUser))
			c.Locals(constants.MiddlewareError, constants.ErrJWTExpired)
			m.Logger.Error("JWT error during authentication in join", zap.Error(err))
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}

		m.Logger.Error("error while checking user identity", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthenticated)
	}

	c.Locals(constants.ContextUid, claims.Subject())
	return c.Next()
}

func (m *Middleware) KratosAuthenticated(c *fiber.Ctx) error {
	kratosId := c.Cookies("ory_kratos_session")
	if kratosId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrKratosIDEmpty)
	}

	kratosUser := config.KratosUserDetails{}
	kratosClient := resty.New().SetBaseURL(m.Config.Kratos.BaseUrl+"/sessions").SetHeader("Cookie", fmt.Sprintf("%v=%v", constants.KratosCookie, kratosId)).SetHeader("accept", "application/json").SetHeader("withCredentials", "true").SetHeader("credentials", "include")

	res, err := kratosClient.R().SetResult(&kratosUser).Get("/whoami")
	if err != nil || res.StatusCode() != http.StatusOK {
		m.Logger.Debug("unauthenticated registration", zap.Any("response from kratos", res.RawResponse))
		return utils.JSONError(c, res.StatusCode(), constants.ErrKratosAuth)
	} else {
		user := models.User{
			KratosID: kratosUser.Identity.ID,
			FirstName: kratosUser.Identity.Traits.Name.First,
			LastName: kratosUser.Identity.Traits.Name.Last,
			Email: kratosUser.Identity.Traits.Email,
		}

		userModel, err := models.InitUserModel(m.Db)
		if err != nil {
			m.Logger.Debug("error while initializing user model", zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
		}

		localUser, err := userModel.GetUserByKratosID(kratosUser.Identity.ID)
		if err != nil {
			m.Logger.Debug("error while getting user information from database while kratos authentication", zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
		}

		user.ID = localUser.ID
		user.Roles = localUser.Roles
		user.Username = localUser.Username

		c.Locals(constants.ContextUser, user)
		c.Locals(constants.ContextUid, user.ID)
		return c.Next()
	}
}
