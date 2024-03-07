package middlewares

import (
	"database/sql"
	"fmt"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	quizHelper "github.com/Improwised/quizz-app/api/helpers/quiz"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PlayedQuizValidationMiddleware struct {
	Config  config.AppConfig
	Logger  *zap.Logger
	Db      *goqu.Database
	helpers *quizHelper.HelperGroup
}

func NewPlayedQuizValidationMiddleware(cfg config.AppConfig, logger *zap.Logger, db *goqu.Database, helper *quizHelper.HelperGroup) PlayedQuizValidationMiddleware {
	return PlayedQuizValidationMiddleware{
		Config:  cfg,
		Logger:  logger,
		Db:      db,
		helpers: helper,
	}
}

func (m *PlayedQuizValidationMiddleware) PlayedQuizValidation(c *fiber.Ctx) error {
	if c.Locals(constants.MiddlewareError) != nil {
		return c.Next()
	}

	invitationCode := quizUtilsHelper.GetString(c.Locals(constants.QuizSessionInvitationCode))

	session, err := m.helpers.ActiveQuizModel.GetSessionByCode(invitationCode)

	if err != nil {
		c.Locals(constants.MiddlewareError, constants.ErrInvitationCodeNotFound)
		m.Logger.Error("error in invitation code", zap.Error(err))
		return c.Next()
	}
	c.Locals(constants.SessionObj, session)
	fmt.Println("hello", session)

	if !session.IsActive || session.ActivatedTo.Valid {

		c.Locals(constants.MiddlewareError, constants.ErrInvitationCodeNotFound)
		m.Logger.Error("invitation code is un-active", zap.Error(fmt.Errorf(constants.ErrInvitationCodeNotFound)))
		return c.Next()
	}

	// get or create user session
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	// if current user is admin of the quiz then no need to create user-played-quiz record
	if userId == session.AdminID {
		return c.Next()
	}

	var userPlayedQuizId uuid.UUID

	if userId == "<nil>" { // anonymous user
		userPlayedQuizId, err = m.helpers.UserPlayedQuizModel.CreateUserPlayedQuiz(sql.NullString{}, session.ID, false)
	} else {
		userPlayedQuizId, _, err = m.helpers.UserPlayedQuizModel.CreateUserPlayedQuizIfNotExists(userId, session.ID)
	}

	if err != nil {
		c.Locals(constants.MiddlewareError, constants.ErrUserQuizSessionValidation)
		m.Logger.Error("Username not provided", zap.Error(err))
		return c.Next()
	}
	c.Locals(constants.CurrentUserQuiz, userPlayedQuizId.String())

	c.Cookie(CreateStrictCookie(constants.CurrentUserQuiz, userPlayedQuizId.String()))
	return c.Next()
}
