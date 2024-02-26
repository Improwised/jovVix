package middlewares

import (
	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type Middleware struct {
	config      config.AppConfig
	logger      *zap.Logger
	db          *goqu.Database
	UserService *services.UserService
}

func NewMiddleware(cfg config.AppConfig, logger *zap.Logger, db *goqu.Database, userService *services.UserService) Middleware {
	return Middleware{
		config:      cfg,
		logger:      logger,
		db:          db,
		UserService: userService,
	}
}
