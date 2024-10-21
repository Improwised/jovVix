package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/pkg/jwt"
	"github.com/Improwised/quizz-app/api/utils"
	resty "github.com/go-resty/resty/v2"
	fiber "github.com/gofiber/fiber/v2"
	j "github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

func (m *Middleware) Authenticated(c *fiber.Ctx) error {
	token := c.Cookies(constants.CookieUser, "")
	kratosToken := c.Cookies(constants.KratosCookie, "")

	if kratosToken == "" {
		if token == "" {
			m.Logger.Debug("returning unauthorized after found empty token")
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}
	} else {
		return m.KratosAuthenticated(c)
	}

	claims, err := jwt.ParseToken(m.Config, token)
	if err != nil {
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			c.Cookie(RemoveCookie(constants.CookieUser))
			m.Logger.Error("JWT error during authentication in join", zap.Error(err))
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}

		m.Logger.Error("error while checking user identity", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthenticated)
	}

	// Check if expiration time is less than 10 minutes than refesh token
	if time.Until(claims.Expiration()) < 10*time.Minute {
		cookieExpirationTime, err := time.ParseDuration(m.Config.Kratos.CookieExpirationTime)
		if err != nil {
			m.Logger.Error("unable to parse the duration for the cookie expiration", zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrKratosCookieTime)
		}

		// Generate a new token
		newToken, err := jwt.CreateToken(m.Config, claims.Subject(), time.Now().Add(time.Hour*2))
		if err != nil {
			m.Logger.Error("error while refreshing token", zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthenticated)
		}

		// Set Coockie
		c.Cookie(&fiber.Cookie{
			Name:    constants.CookieUser,
			Value:   newToken,
			Expires: time.Now().Add(cookieExpirationTime),
		})
	}

	c.Locals(constants.ContextUid, claims.Subject())
	return c.Next()
}

func (m *Middleware) KratosAuthenticated(c *fiber.Ctx) error {
	kratosId := c.Cookies("ory_kratos_session")
	if kratosId == "" {
		return utils.JSONError(c, http.StatusUnauthorized, constants.ErrKratosIDEmpty)
	}

	kratosUser := config.KratosUserDetails{}
	kratosClient := resty.New().SetBaseURL(m.Config.Kratos.BaseUrl+"/sessions").SetHeader("Cookie", fmt.Sprintf("%v=%v", constants.KratosCookie, kratosId)).SetHeader("accept", "application/json").SetHeader("withCredentials", "true").SetHeader("credentials", "include")

	res, err := kratosClient.R().SetResult(&kratosUser).Get("/whoami")
	if err != nil || res.StatusCode() != http.StatusOK {
		m.Logger.Debug("unauthenticated registration", zap.Any("response from kratos", res.RawResponse))
		return utils.JSONError(c, res.StatusCode(), constants.ErrKratosAuth)
	}

	m.Logger.Debug("userModel.GetUserByKratosID called", zap.Any("kratosID", kratosUser.Identity.ID))
	user, err := m.userModel.GetUserByKratosID(kratosUser.Identity.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			m.Logger.Error(constants.ErrGetUser, zap.Error(err))
			return utils.JSONError(c, http.StatusNotFound, constants.ErrGetUser)
		}
		m.Logger.Error(constants.ErrGetUser, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetUser)
	}
	m.Logger.Debug("userModel.GetUserByKratosID success", zap.Any("user", user))

	c.Locals(constants.ContextUser, user)
	c.Locals(constants.ContextUid, user.ID)
	c.Locals(constants.KratosID, kratosUser.Identity.ID)
	return c.Next()

}
