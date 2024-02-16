package v1

import (
	"fmt"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

func CreateUser(db *goqu.Database, logger *zap.Logger, userObj models.User, retrying bool) (models.User, error) {
	userModel, err := models.InitUserModel(db)

	if err != nil {
		return userObj, err
	}

	isUnique, err := userModel.IsUniqueEmail(userObj.Email)

	if err != nil {
		return userObj, fmt.Errorf("someError occurred during register user %v", err)
	}

	if !isUnique {
		return userObj, fmt.Errorf("email is already registered")
	}

	userSvc := services.NewUserService(&userModel)

	copyUserObj, err := userSvc.RegisterUser(userObj, events.NewEventBus(logger))
	
	if err != nil {
		
		// Check is there unique-key error
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			
			if !(retrying && pqErr.Constraint == constants.UserUkey) {
				return userObj, fmt.Errorf("username (%s) already registered", userObj.Username)
			}
			
			copyUserObj.Password = userObj.Password

			copyUserObj.Username = utils.GenerateNewStringHavingSuffixName(userObj.Username, 5, 12)

			copyUserObj, err = userSvc.RegisterUser(copyUserObj, events.NewEventBus(logger))

			if err != nil {
				return userObj, fmt.Errorf("SomeError during register admin with new username %s", userObj.Username)
			}

		}

	}

	userObj.ID = copyUserObj.ID
	userObj.Username = copyUserObj.Username

	return userObj, err
}
