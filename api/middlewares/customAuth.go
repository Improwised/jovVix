package middlewares

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Improwised/quizz-app/api/constants"
	quizController "github.com/Improwised/quizz-app/api/controllers/api/v1"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/jwt"
	"github.com/Improwised/quizz-app/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	j "github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

// path 1. If user set in cookie
// check if jwk
//
//			A. exists:- verify jwk
//				a. jwk ok:- get user
//								1. found:- set role in context and session storage
//	                         2. not found :- TODO: send error
//				b. jwk not ok:- fail message return/ login again or join as user
//	  	B. not exists :- trying to get userName from query and create new_user get userID and role and set in context and cookie
//
// path 2. If user not set in cookie
//
//	B. not exists :- trying to get userName from query and create new_user get userID and role and set in context and cookie
func (m *Middleware) AuthenticationAndRoleAssignment(c *fiber.Ctx) error {

	if m.config.Kratos.IsEnabled {
		sessionID := c.Cookies("ory_kratos_session")
		if sessionID == "" {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}
		c.Locals(constants.KratosID, sessionID)
		return c.Next()
	}

	token := c.Cookies(constants.CookieUser, "")
	var err error

	if token != "" {
		err = AuthHavingTokenHandler(m, c, token)
	} else {
		err = AuthHavingNotTokenHandler(m, c)
	}

	if err != nil {
		return err
	}

	return c.Next()
}

func AuthHavingTokenHandler(m *Middleware, c *fiber.Ctx, token string) error {
	// JWK verification
	claims, err := jwt.ParseToken(m.config, token)

	if err != nil {
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		}

		m.logger.Error("error while checking user identity", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUnauthenticated)
	}
	c.Locals(constants.ContextUid, claims.Subject())

	userObj, err := m.userService.GetUser(c.Locals(constants.ContextUid).(string))

	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.InvalidCredentials)
		}
		return utils.JSONFail(c, http.StatusInternalServerError, err)
	}

	// Set user config as current session

	c.Cookie(CreateStrictCookie(constants.CookieUser, userObj.Username))
	c.Locals(constants.ContextUid, userObj.ID)
	c.Locals(constants.ContextUser, userObj)

	return nil
}

func AuthHavingNotTokenHandler(m *Middleware, c *fiber.Ctx) error {
	// get userName from query
	userName := c.Query(constants.UserName, "")

	if userName == "" {
		return utils.JSONFail(c, http.StatusBadRequest, constants.UsernameRequired)
	}

	userObj := models.User{
		Username: userName,
		Roles:    "user",
	}
	userObj, err := quizController.CreateQuickUser(m.db, m.logger, userObj, true, false)

	if err != nil {
		return utils.JSONFail(c, http.StatusInternalServerError, err)
	}

	c.Cookie(CreateStrictCookie(constants.CookieUser, userObj.Username))
	c.Locals(constants.ContextUid, userObj.ID)
	c.Locals(constants.ContextUser, userObj)

	return nil
}

func CreateStrictCookie(key, value string) *fiber.Cookie {
	cookieUserName := new(fiber.Cookie)
	cookieUserName.Name = key
	cookieUserName.Value = value
	cookieUserName.HTTPOnly = true
	cookieUserName.SessionOnly = true
	cookieUserName.SameSite = "Strict"
	cookieUserName.Expires = time.Now().Add(2 * time.Hour)
	return cookieUserName
}
