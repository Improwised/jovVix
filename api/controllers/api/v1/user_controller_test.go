package v1_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGuestUser(t *testing.T) {
	username := "username"
	avatarName := "Chase"
	t.Run("create user without giving avatar", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Post(fmt.Sprintf("/api/v1/user/%s", username))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})

	t.Run("create user with valid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Post(fmt.Sprintf("/api/v1/user/%s?avatar_name=%s", username, avatarName))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})

	t.Cleanup(func() {
		_, err := db.Exec("delete from users where username='username'")
		assert.Nil(t, err)
	})
}

func TestGetUserMeta(t *testing.T) {

	t.Run("create user with valid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/user/who")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}
