package middlewares

import (
	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/models"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type Middleware struct {
	Config             config.AppConfig
	Logger             *zap.Logger
	Db                 *goqu.Database
	userModel          *models.UserModel
	sharedQuizzesModel *models.SharedQuizzesModel
}

func NewMiddleware(cfg config.AppConfig, logger *zap.Logger, db *goqu.Database) Middleware {

	userModel, _ := models.InitUserModel(db, logger)
	sharedQuizzesModel := models.InitSharedQuizzesModel(db, logger)

	return Middleware{
		Config:             cfg,
		Logger:             logger,
		Db:                 db,
		userModel:          &userModel,
		sharedQuizzesModel: sharedQuizzesModel,
	}
}
