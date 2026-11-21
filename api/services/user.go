package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
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

	if err := deleteKratosIdentity(userSvc.config.Kratos.AdminUrl, kratosId); err != nil {
		userSvc.logger.Error("error while deleting identity from kratos", zap.Error(err))
		return fmt.Errorf("failed to delete user from kratos for %s", kratosId)
	}

	isOk = true

	return nil
}

// deleteKratosIdentity deletes the identity through the Kratos admin API.
// adminUrl is SERVE_ADMIN_BASE_URL and includes the /admin prefix, as in api/cli/delete-orphaned.go.
// Redirects are not followed: the admin API has no authentication, so a misconfigured URL must fail
// instead of being sent somewhere else.
func deleteKratosIdentity(adminUrl, kratosId string) error {
	if strings.TrimSpace(adminUrl) == "" {
		return fmt.Errorf("SERVE_ADMIN_BASE_URL is not set")
	}

	res, err := resty.New().
		SetBaseURL(strings.TrimRight(adminUrl, "/")+"/identities").
		SetHeader("accept", "application/json").
		SetRedirectPolicy(resty.NoRedirectPolicy()).
		R().
		SetPathParam("id", kratosId).
		Delete("/{id}")
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("kratos admin API returned %d", res.StatusCode())
	}
	return nil
}
