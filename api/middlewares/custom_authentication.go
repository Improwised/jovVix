package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	quizController "github.com/Improwised/quizz-app/api/controllers/api/v1"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
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
//	                         2. not found :- send error
//				b. jwk not ok:- fail message return/ login again or join as user
//	  	B. not exists :- trying to get userName from query and create new_user get userID and role and set in context and cookie
//
// path 2. If user not set in cookie
//
//	B. not exists :- trying to get userName from query and create new_user get userID and role and set in context and cookie
func (m *Middleware) CustomAuthenticated(c *fiber.Ctx) error {

	if c.Locals(constants.MiddlewareError) != nil {
		return c.Next()
	}

	token := c.Cookies(constants.CookieUser, "cookie not available")

	if token == "cookie not available" {
		AuthHavingNoTokenHandler(m, c)
	} else {
		AuthHavingTokenHandler(m, c, token)
	}

	return c.Next()
}

func (m *Middleware) CustomAdminAuthenticated(c *fiber.Ctx) error {

	if c.Locals(constants.MiddlewareError) != nil {
		return c.Next()
	}

	token := c.Cookies(constants.CookieUser, "cookie not available")

	if token != "cookie not available" {
		AuthHavingTokenHandler(m, c, token)
	} else {
		c.Locals(constants.MiddlewareError, constants.Unauthenticated)
	}

	return c.Next()
}

func (m *Middleware) CheckSessionId(c *fiber.Ctx) error {

	if c.Locals(constants.MiddlewareError) != nil {
		return c.Next()
	}

	// get session id from param
	sessionId := c.Params(constants.SessionIDParam)

	c.Locals(constants.SessionIDParam, sessionId)
	return c.Next()
}

func (m *Middleware) CheckSessionCode(c *fiber.Ctx) error {

	if c.Locals(constants.MiddlewareError) != nil {
		return c.Next()
	}

	// get session code from param
	code := c.Params(constants.QuizSessionInvitationCode)

	if !quizUtilsHelper.IsValidCode(code) {
		c.Locals(constants.MiddlewareError, constants.ErrInvitationCodeInWrongFormat)
		return c.Next()
	}

	c.Locals(constants.QuizSessionInvitationCode, code)
	return c.Next()
}

func AuthHavingTokenHandler(m *Middleware, c *fiber.Ctx, token string) {
	// JWK verification
	claims, err := jwt.ParseToken(m.Config, token)

	if err != nil {
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			c.Cookie(RemoveCookie(constants.CookieUser))
			c.Locals(constants.MiddlewareError, constants.ErrJWTExpired)
			m.Logger.Error("JWT error during authentication in join", zap.Error(err))
			return
		}

		m.Logger.Error("Error while checking user identity in join", zap.Error(err))
		c.Locals(constants.MiddlewareError, constants.ErrUnauthenticated)
		return
	}

	c.Locals(constants.ContextUid, claims.Subject())

	userObj, err := m.UserService.GetUser(claims.Subject())

	if err != nil {
		if err == sql.ErrNoRows {
			c.Locals(constants.MiddlewareError, constants.InvalidCredentials)
			m.Logger.Error(fmt.Sprintf("User not found, userID %s", claims.Subject()), zap.Error(err))
			return
		}

		c.Locals(constants.MiddlewareError, constants.UnknownError)
		m.Logger.Error(fmt.Sprintf("Unknown DB error, userID %s", claims.Subject()), zap.Error(err))
		return
	}

	// c.Cookie(CreateStrictCookie(constants.CookieUser, token))
	c.Locals(constants.ContextUid, userObj.ID)
	c.Locals(constants.ContextUser, userObj)
}

func AuthHavingNoTokenHandler(m *Middleware, c *fiber.Ctx) {
	// get userName from query
	userName := c.Query(constants.UserName, "")

	if userName == "" {
		c.Locals(constants.MiddlewareError, constants.UsernameRequired)
		m.Logger.Error("Username not provided", zap.Error(fmt.Errorf(constants.UsernameRequired)))
		return
	}

	userObj := models.User{
		Username: userName,
		Roles:    "user",
	}

	userObj, err := quizController.CreateQuickUser(m.Db, m.Logger, userObj, true, false)

	if err != nil {
		c.Locals(constants.MiddlewareError, err)
		m.Logger.Error(fmt.Sprintf("Error in register user %v", userObj), zap.Error(err))
		return
	}

	err = CreateNewUserToken(c, m.Config, userObj, m.Logger)
	if err != nil {
		c.Locals(constants.MiddlewareError, err)
		m.Logger.Error(fmt.Sprintf("Error in register user %v", userObj), zap.Error(err))
		return
	}
}

func CreateStrictCookie(key, value string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.HTTPOnly = true
	cookie.SessionOnly = true
	if os.Getenv("IS_DEVELOPMENT") == "true" {
		cookie.Secure = false
	} else {
		cookie.Secure = true
	}
	cookie.SameSite = "Strict"
	return cookie
}

func RemoveCookie(key string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.Value = "" // Or set a generic value like "deleted"
	cookie.Expires = time.Now().Add(-1 * time.Second)
	return cookie
}

func CreateNewUserToken(c *fiber.Ctx, cfg config.AppConfig, user models.User, logger *zap.Logger) error {
	token, err := jwt.CreateToken(cfg, user.ID, time.Now().Add(time.Hour*2))
	if err != nil {
		logger.Error("error while creating token", zap.Error(err), zap.Any("id", user.ID))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrLoginUser)
	}
	c.Cookie(CreateStrictCookie(constants.CookieUser, token))
	c.Locals(constants.ContextUid, user.ID)
	c.Locals(constants.ContextUser, user)
	return nil
}
