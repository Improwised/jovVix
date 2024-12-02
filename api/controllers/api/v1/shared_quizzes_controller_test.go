package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Improwised/quizz-app/api/utils"
	"github.com/stretchr/testify/assert"
)

var shareQuizId string

func shareQuiz(t *testing.T, email, permission string) {
	req := map[string]string{
		"email":      email,
		"permission": permission,
	}

	res, err := client.
		R().
		EnableTrace().
		SetBody(req).
		Post(fmt.Sprintf("/api/v1/shared_quizzes/%s", quizId))

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode())

	var sharequizResponse ResponseQuiz
	err = json.Unmarshal(res.Body(), &sharequizResponse)
	assert.Nil(t, err, "Error in parsing response body")

	shareQuizId = sharequizResponse.Data
}

func TestShareQuiz(t *testing.T) {
	invalidQuizId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"
	t.Run("share quiz with invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Post(fmt.Sprintf("/api/v1/shared_quizzes/%s", invalidQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("share quiz with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		shareQuiz(t, "demosharequizuser@example.com", "read")
	})
}

func TestListQuizAuthorizedUsers(t *testing.T) {
	invalidQuizId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"
	t.Run("list authorized quiz user with invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/shared_quizzes/%s", invalidQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("list authorized quiz user with valid input", func(t *testing.T) {
		// share quiz to users
		createQuiz(t, quizTitle)
		shareQuiz(t, "demosharequizuser@example.com", "read")
		shareQuiz(t, "demosharequizuser2@example.com", "share")

		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/shared_quizzes/%s", quizId))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var listQuizAuthorizedUsers utils.ResponseListQuizAuthorizedUsers
		err = json.Unmarshal(res.Body(), &listQuizAuthorizedUsers.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.GreaterOrEqual(t, len(listQuizAuthorizedUsers.Body.Data), 2, res)
	})
}

func TestUpdateUserPermissionOfQuiz(t *testing.T) {
	invalidQuizId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"
	t.Run("update quiz permission with invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Put(fmt.Sprintf("/api/v1/shared_quizzes/%s", invalidQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("update quiz permission with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		shareQuiz(t, "demosharequizuser@example.com", "read")

		params := map[string]string{
			"shared_quiz_id": shareQuizId,
		}

		req := map[string]interface{}{
			"email":      "demosharequizuser@example.com",
			"permission": "share",
		}

		res, err := client.
			R().
			EnableTrace().
			SetQueryParams(params).
			SetBody(req).
			Put(fmt.Sprintf("/api/v1/shared_quizzes/%s", quizId))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}

func TestDeleteUserPermissionOfQuiz(t *testing.T) {
	invalidQuizId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"
	t.Run("delete quiz permission with invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Delete(fmt.Sprintf("/api/v1/shared_quizzes/%s", invalidQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("delete quiz permission with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		shareQuiz(t, "demosharequizuser@example.com", "read")

		params := map[string]string{
			"shared_quiz_id": shareQuizId,
		}

		res, err := client.
			R().
			EnableTrace().
			SetQueryParams(params).
			Delete(fmt.Sprintf("/api/v1/shared_quizzes/%s", quizId))
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}

func TestListSharedQuizzes(t *testing.T) {
	t.Run("list shared quizzes of user with invalid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		shareQuiz(t, "demosharequizuser@example.com", "read")

		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/shared_quizzes")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("list shared quizzes of user with valid input (type=shared_with_me)", func(t *testing.T) {
		params := map[string]string{
			"type": "shared_with_me",
		}

		res, err := client.
			R().
			EnableTrace().
			SetQueryParams(params).
			Get("/api/v1/shared_quizzes")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var listSharedQuizzes utils.ResponseListSharedQuizzes
		err = json.Unmarshal(res.Body(), &listSharedQuizzes.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, 0, len(listSharedQuizzes.Body.Data), res)
	})

	t.Run("list shared quizzes of user with valid input (type=shared_by_me)", func(t *testing.T) {
		// share quiz to users
		createQuiz(t, quizTitle)
		shareQuiz(t, "demosharequizuser@example.com", "read")
		createQuiz(t, quizTitle)
		shareQuiz(t, "demosharequizuser@example.com", "share")

		params := map[string]string{
			"type": "shared_by_me",
		}

		res, err := client.
			R().
			EnableTrace().
			SetQueryParams(params).
			Get("/api/v1/shared_quizzes")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var listSharedQuizzes utils.ResponseListSharedQuizzes
		err = json.Unmarshal(res.Body(), &listSharedQuizzes.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.GreaterOrEqual(t, len(listSharedQuizzes.Body.Data), 2, res)
	})
}
