package v1_test

import (
	"net/http"
	"testing"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/Improwised/quizz-app/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/assert"
)

func TestDoAuth(t *testing.T) {

	cfg := config.LoadTestEnv()

	db, err := database.Connect(cfg.DB)
	assert.Nil(t, err)

	t.Run("Do authentication when user logged in", func(t *testing.T) {
		var actual utils.ResponseAuthnUser
		_, err := db.Insert("users").Rows(
			goqu.Record{
				"id":         "coq5km6bcbvvgbgfuek0",
				"first_name": "admin",
				"last_name":  "xyz",
				"email":      "adminxyz@gmail.com",
				"password":   "RZo5(uXD<3#aH0",
				"roles":      "admin",
				"username":   "adminxyz123",
			},
		).Executor().Exec()

		assert.Nil(t, err)

		req := structs.ReqLoginUser{
			Email:    "adminxyz@gmail.com",
			Password: "RZo5(uXD<3#aH0",
		}

		expected := models.User{
			ID:        "coq5km6bcbvvgbgfuek0",
			KratosID:  "",
			FirstName: "admin",
			LastName:  "xyz",
			Email:     "adminxyz@gmail.com",
			Username:  "",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			SetResult(&actual.Body).
			Post("/api/v1/login")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
		assert.Equal(t, expected.ID, actual.Body.Data.ID, "expected and actual ID is not equal")
		assert.Equal(t, expected.KratosID, "", "expected and actual KratosID is not equal")
		assert.Equal(t, expected.FirstName, actual.Body.Data.FirstName, "expected and actual FirstName is not equal")
		assert.Equal(t, expected.LastName, actual.Body.Data.LastName, "expected and actual LastName is not equal")
		assert.Equal(t, expected.Email, actual.Body.Data.Email, "expected and actual Email is not equal")
		assert.Equal(t, expected.Username, "", "expected and actual Username is not equal")
		cookies := res.Cookies()
		assert.NotEqual(t, len(cookies), 0)
	})
}
