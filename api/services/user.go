package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/doug-martin/goqu/v9"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type UserService struct {
	userPlayedQuizModel *models.UserPlayedQuizModel
	quizModel           *models.QuizModel
	activeQuizModel     *models.ActiveQuizModel
	userModel           *models.UserModel
	sharedQuizzesModel  *models.SharedQuizzesModel
	db                  *goqu.Database
	logger              *zap.Logger
	config              config.AppConfig
}

func NewUserService(db *goqu.Database, logger *zap.Logger, config config.AppConfig) (*UserService, error) {
	userPlayedQuizModel := models.InitUserPlayedQuizModel(db)
	quizModel := models.InitQuizModel(db)
	activeQuizModel := models.InitActiveQuizModel(db, logger)
	sharedQuizzesModel := models.InitSharedQuizzesModel(db, logger)
	userModel, err := models.InitUserModel(db, logger)
	if err != nil {
		return nil, err
	}

	return &UserService{
		userPlayedQuizModel: userPlayedQuizModel,
		quizModel:           quizModel,
		activeQuizModel:     activeQuizModel,
		sharedQuizzesModel:  sharedQuizzesModel,
		userModel:           &userModel,
		db:                  db,
		logger:              logger,
		config:              config,
	}, nil
}

// This function is delete all user data from database
func (userSvc *UserService) DeleteUserDataById(userId, kratosId string) error {
	isOk := false
	transaction, err := userSvc.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				userSvc.logger.Error("error during commit in delete user", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				userSvc.logger.Error("error during rollback in delete user", zap.Error(err))
			}
		}
	}()

	// Delete shared quizzes and shared with me quizzes also
	err = userSvc.sharedQuizzesModel.RemoveSharedQuizPermissionsByUserId(transaction, userId)
	if err != nil {
		return err
	}

	// Delete played quizzes
	err = userSvc.userPlayedQuizModel.DeleteUserPlayedQuizzesAndReponseByUserId(transaction, userId)
	if err != nil {
		return err
	}

	// Delete active quizzes
	err = userSvc.activeQuizModel.DeleteActiveQuizzesAndRelatedDataByUserId(transaction, userId)
	if err != nil {
		return err
	}

	// Delete created quizzes
	err = userSvc.quizModel.DeleteCreatedQuizzesByUserId(transaction, userId)
	if err != nil {
		return err
	}

	// Delete user from users table
	deletedkratosId, err := userSvc.userModel.DeleteUserById(transaction, userId)
	if err != nil {
		return err
	}

	if strings.TrimSpace(deletedkratosId) != kratosId {
		userSvc.logger.Error(constants.Unauthenticated)
		return fmt.Errorf(constants.Unauthenticated)
	}

	// delete user data from kratos database
	kratosClient := resty.New().SetBaseURL(userSvc.config.Kratos.BaseUrl+"/admin/identities").SetHeader("accept", "application/json")

	res, err := kratosClient.R().Delete(fmt.Sprintf("/%s", kratosId))
	if err != nil || res.StatusCode() != http.StatusNoContent {
		userSvc.logger.Error("unauthenticated registration", zap.Any("response from kratos", res.RawResponse), zap.Error(err), zap.Any("kratos response", res))
		return fmt.Errorf("failed to delete user from kratos for %s", kratosId)
	}

	isOk = true

	return nil
}
