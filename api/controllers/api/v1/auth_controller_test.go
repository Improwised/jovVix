package v1_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Improwised/quizz-app/api/utils"
	"github.com/stretchr/testify/assert"
)

func TestDoKratosAuth(t *testing.T) {

	t.Run("get user with valid input", func(t *testing.T) {
		client.
			R().
			EnableTrace().
			Get("/api/v1/kratos/auth")
	})
}

func TestGetRegisteredUser(t *testing.T) {

	t.Run("get user with valid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/kratos/whoami")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var user utils.ResponseGetRegisteredUser
		err = json.Unmarshal(res.Body(), &user.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "exampletestuser@example.com", user.Body.Data.Identity.Traits.Email, res)
		assert.Equal(t, "John", user.Body.Data.Identity.Traits.Name.First, res)
		assert.Equal(t, "Doe", user.Body.Data.Identity.Traits.Name.Last, res)
	})
}

func TestUpadateRegisteredUser(t *testing.T) {
	t.Run("update user details with invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Put("/api/v1/kratos/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("update user details with invalid input (no email is in payload)", func(t *testing.T) {
		req := map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/kratos/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("update user details with invalid input (email is not in formate)", func(t *testing.T) {
		req := map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "demogmail",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/kratos/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("update user details with valid input", func(t *testing.T) {
		req := map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "exampletestuser@example.com",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put("/api/v1/kratos/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}
