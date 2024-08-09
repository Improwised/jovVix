package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
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
//	                         2. not found :- send error
//				b. jwk not ok:- fail message return/ login again or join as user
//	  	B. not exists :- trying to get userName from query and create new_user get userID and role and set in context and cookie
//
// path 2. If user not set in cookie
//
//	B. not exists :- trying to get userName from query and create new_user get userID and role and set in context and cookie
func (m *Middleware) CustomAuthenticated(c *fiber.Ctx) error {

	token := c.Cookies(constants.CookieUser, "cookie not available")
	kratosToken := c.Cookies(constants.KratosCookie, "cookie not available")
	if kratosToken == "cookie not available" {
		if token == "cookie not available" {
			return utils.JSONFail(c, http.StatusUnauthorized, constants.Unauthenticated)
		} else {
			return AuthHavingTokenHandler(m, c, token)
		}
	} else {
		return m.KratosAuthenticated(c)
	}
}

func (m *Middleware) CheckSessionId(c *fiber.Ctx) error {

	// get session id from param
	sessionId := c.Params(constants.SessionIDParam)

	c.Locals(constants.SessionIDParam, sessionId)
	return c.Next()
}

func (m *Middleware) CheckSessionCode(c *fiber.Ctx) error {

	// get session code from param
	code := c.Params(constants.QuizSessionInvitationCode)
	if !quizUtilsHelper.IsValidCode(code) {
		m.Logger.Error(constants.ErrInvitationCodeInWrongFormat)
		return utils.JSONFail(c, http.StatusUnauthorized, constants.ErrInvitationCodeInWrongFormat)
	}

	c.Locals(constants.QuizSessionInvitationCode, code)
	return c.Next()
}

func AuthHavingTokenHandler(m *Middleware, c *fiber.Ctx, token string) error {
	// JWK verification
	claims, err := jwt.ParseToken(m.Config, token)
	if err != nil {
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			c.Cookie(RemoveCookie(constants.CookieUser))
			m.Logger.Error(constants.ErrJWTExpired, zap.Error(err))
			return utils.JSONFail(c, http.StatusUnauthorized, constants.ErrJWTExpired)
		}
		m.Logger.Error("Error while checking user identity in join", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, "Error while checking user identity in join")
	}

	userObj, err := m.UserService.GetUser(claims.Subject())
	if err != nil {
		if err == sql.ErrNoRows {
			m.Logger.Error(fmt.Sprintf("User not found, userID %s", claims.Subject()), zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, constants.InvalidCredentials)
		}
		m.Logger.Error(fmt.Sprintf("Unknown DB error, userID %s", claims.Subject()), zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrGetUser)
	}

	c.Locals(constants.ContextUid, userObj.ID)
	c.Locals(constants.ContextUser, userObj)

	return c.Next()
}

func RemoveCookie(key string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.Value = "" // Or set a generic value like "deleted"
	cookie.Expires = time.Now().Add(-1 * time.Second)
	return cookie
}
