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
		fmt.Println(err)
		return userObj, fmt.Errorf("SomeError occurred during register user")
	} else if isUnique {
		return userObj, fmt.Errorf("Unique key violation on Email")
	}

	userSvc := services.NewUserService(&userModel)

	userObj, err = userSvc.RegisterUser(userObj, events.NewEventBus(logger))

	if err != nil {

		// Check is there unique-key error
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {

			if !(retrying && pqErr.Constraint == constants.UserUkey) {
				return userObj, fmt.Errorf("Unique key violation on %v", pqErr.Constraint)
			}

			// force flag is set then it tries agin with manipulated username trying to make new admin-username

			userObj.Username = utils.GenerateNewStringHavingSuffixName(userObj.Username, 5, 12)

			_, err = userSvc.RegisterUser(userObj, events.NewEventBus(logger))

			if err != nil {
				return userObj, fmt.Errorf("SomeError during register admin with new username %s", userObj.Username)
			}
		}

	}

	return userObj, err
}
