package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Improwised/quizz-app/api/config"
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

	c.Locals(constants.MiddlewareFail, nil)

	token := c.Cookies(constants.CookieUser, "cookie not available")

	if token != "cookie not available" {
		AuthHavingTokenHandler(m, c, token)
	} else {
		AuthHavingNoTokenHandler(m, c)
	}

	return c.Next()
}

func AuthHavingTokenHandler(m *Middleware, c *fiber.Ctx, token string) {
	// JWK verification
	claims, err := jwt.ParseToken(m.config, token)

	if err != nil {
		if errors.Is(err, j.ErrInvalidJWT()) || errors.Is(err, j.ErrTokenExpired()) {
			c.Cookie(RemoveUserToken(constants.CookieUser))
			c.Locals(constants.MiddlewareFail, constants.Unauthenticated)
			m.logger.Error("JWT error during authentication in join", zap.Error(err))
			return
		}

		m.logger.Error("Error while checking user identity in join", zap.Error(err))
		c.Locals(constants.MiddlewareFail, constants.ErrUnauthenticated)
		return
	}

	c.Locals(constants.ContextUid, claims.Subject())

	userObj, err := m.userService.GetUser(claims.Subject())

	if err != nil {
		if err == sql.ErrNoRows {
			c.Locals(constants.MiddlewareFail, constants.InvalidCredentials)
			m.logger.Error(fmt.Sprintf("User not found, userID %s", claims.Subject()), zap.Error(err))
			return
		}

		c.Locals(constants.MiddlewareFail, constants.UnknownError)
		m.logger.Error(fmt.Sprintf("Unknown DB error, userID %s", claims.Subject()), zap.Error(err))
		return
	}

	// Set user config as current session

	c.Cookie(CreateStrictCookie(constants.CookieUser, token))
	c.Locals(constants.ContextUid, userObj.ID)
	c.Locals(constants.ContextUser, userObj)
}

func AuthHavingNoTokenHandler(m *Middleware, c *fiber.Ctx) {
	// get userName from query
	userName := c.Query(constants.UserName, "")

	if userName == "" {
		c.Locals(constants.MiddlewareFail, constants.UsernameRequired)
		m.logger.Error("Username not provided", zap.Error(fmt.Errorf(constants.UsernameRequired)))
		return
	}

	userObj := models.User{
		Username: userName,
		Roles:    "user",
	}

	userObj, err := quizController.CreateQuickUser(m.db, m.logger, userObj, true, false)

	if err != nil {
		c.Locals(constants.MiddlewareFail, err)
		m.logger.Error(fmt.Sprintf("Error in register user %v", userObj), zap.Error(err))

		return
	}

	err = CreateNewUserToken(c, m.config, userObj, m.logger)
	if err != nil {
		c.Locals(constants.MiddlewareFail, err)
		m.logger.Error(fmt.Sprintf("Error in register user %v", userObj), zap.Error(err))
		return
	}
}

func CreateStrictCookie(key, value string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.HTTPOnly = true
	cookie.SessionOnly = true
	cookie.Secure = true
	cookie.SameSite = "Strict"
	cookie.Expires = time.Now().Add(2 * time.Hour)
	return cookie
}

func RemoveUserToken(key string) *fiber.Cookie {
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

// func MiddlewareFail(c *websocket.Conn, refStr string, response string, other any) error {
// 	utils.WsJSONFail(c, "Authentication", refStr, response, other)
// 	time.Sleep(100)
// 	return authError
// }
