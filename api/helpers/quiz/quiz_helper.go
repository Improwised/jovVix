package quizHelper

import (
	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/doug-martin/goqu/v9"
)

type HelperStructs struct {
	UserService      *services.UserService
	RoleModel        *models.RoleModel
	QuizModel        *models.QuizModel
	QuizSessionModel *models.QuizSessionModel
	QuestionModel    *models.QuestionAnswerModel
	PubSubModel      *models.PubSubModel
}

func InitHelper(db *goqu.Database, pubSubCfg config.RedisClientConfig) (*HelperStructs, error) {
	userModel, err := models.InitUserModel(db)

	if err != nil {
		return nil, err
	}

	userService := services.NewUserService(&userModel)

	sessionModel, err := models.InitQuizSessionModel(db)
	if err != nil {
		return nil, err
	}

	roleModel := models.InitRoleModel(db)
	quizModel := models.InitQuizModel(db)
	questionModel := models.InitQuestionAnswerModel(db)

	pubSubClientModel, err := models.InitPubSubModel(pubSubCfg.RedisAddr+":"+pubSubCfg.RedisPort, pubSubCfg.RedisPass, pubSubCfg.RedisDb)

	if err != nil {
		return nil, err
	}

	return &HelperStructs{
		UserService:      userService,
		RoleModel:        roleModel,
		QuizModel:        quizModel,
		QuizSessionModel: sessionModel,
		QuestionModel:    questionModel,
		PubSubModel:      pubSubClientModel,
	}, nil
}
