package middlewares

import (
	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type Middleware struct {
	Config    config.AppConfig
	Logger    *zap.Logger
	Db        *goqu.Database
	userModel *models.UserModel
}

func NewMiddleware(cfg config.AppConfig, logger *zap.Logger, db *goqu.Database) Middleware {

	userModel, _ := models.InitUserModel(db, logger)

	return Middleware{
		Config:    cfg,
		Logger:    logger,
		Db:        db,
		userModel: &userModel,
	}
}
