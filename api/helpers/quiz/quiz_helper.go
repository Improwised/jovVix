package quizHelper

import (
	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

type HelperGroup struct {
	UserService           *services.UserService
	RoleModel             *models.RoleModel
	QuizModel             *models.QuizModel
	ActiveQuizModel       *models.ActiveQuizModel
	UserQuizResponseModel *models.UserQuizResponseModel
	QuestionModel         *models.QuestionModel
	UserPlayedQuizModel   *models.UserPlayedQuizModel
	PubSubModel           *models.PubSubModel
}

func InitHelper(db *goqu.Database, pubSubCfg config.RedisClientConfig, logger *zap.Logger) (*HelperGroup, error) {
	userModel, err := models.InitUserModel(db)

	if err != nil {
		return nil, err
	}

	userService := services.NewUserService(&userModel)

	sessionModel := models.InitActiveQuizModel(db, logger)
	roleModel := models.InitRoleModel(db)
	quizModel := models.InitQuizModel(db)
	userQuizResponseModel := models.InitUserQuizResponseModel(db)
	userPlayedQuizModel := models.InitUserPlayedQuizModel(db)
	questionModel := models.InitQuestionModel(db, logger)

	pubSubClientModel, err := models.InitPubSubModel(pubSubCfg.RedisAddr+":"+pubSubCfg.RedisPort, pubSubCfg.RedisPass, pubSubCfg.RedisDb)

	if err != nil {
		return nil, err
	}

	return &HelperGroup{
		UserService:           userService,
		RoleModel:             roleModel,
		QuizModel:             quizModel,
		ActiveQuizModel:       sessionModel,
		UserQuizResponseModel: userQuizResponseModel,
		QuestionModel:         questionModel,
		PubSubModel:           pubSubClientModel,
		UserPlayedQuizModel:   userPlayedQuizModel,
	}, nil
}
