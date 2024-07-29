package services

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	err := os.Chdir("../")
	assert.Nil(t, err)

	cfg := config.LoadTestEnv()

	db, err := database.Connect(cfg.DB)
	assert.Nil(t, err)

	userModel, err := models.InitUserModel(db)
	assert.Nil(t, err)

	userSvc := NewUserService(&userModel)
	email := "someone11@improwised.com"
	uid := ""

	t.Run("test user register", func(t *testing.T) {
		events := events.NewMockIEvents(t)
		regdEvent := structs.EventUserRegistered{
			Email: email,
		}

		events.EXPECT().Publish(constants.EventUserRegistered, regdEvent)
		data := models.User{
			FirstName: "someone",
			LastName:  "example",
			Email:     email,
			Password: sql.NullString{
				String: "sdhd^72AAAyuZM",
				Valid:  true,
			},
			Roles:    "admin",
			Username: "different",
		}

		resp, err := userSvc.RegisterUser(data, events)
		assert.Nil(t, err)

		assert.Equal(t, data.Email, resp.Email)
		uid = resp.ID
	})

	t.Run("getUser should return user not exists error", func(t *testing.T) {
		_, err := userSvc.GetUser("11")
		assert.Equal(t, sql.ErrNoRows, err)
	})

	t.Run("get exist user by id", func(t *testing.T) {
		user, err := userSvc.GetUser(uid)
		assert.Nil(t, err)
		assert.Equal(t, email, user.Email)
	})

	t.Cleanup(func() {
		_, err = db.Exec("delete from users where email = $1", email)
		assert.Nil(t, err)
	})
}
