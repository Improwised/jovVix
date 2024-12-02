package v1_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Improwised/quizz-app/api/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetScore(t *testing.T) {

	t.Run("get final score for user with invalid input (user_played_quiz query params is not given)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/final_score/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("get final score for user with invalid input (user_played_quiz query params is invalid)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			SetQueryParam("user_played_quiz", invaliduserPlayedQuizId).
			Get("/api/v1/final_score/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseFinalScoreForUser utils.ResponseFinalScoreForUser
		err = json.Unmarshal(res.Body(), &responseFinalScoreForUser.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "success", responseFinalScoreForUser.Body.Status, res)
		assert.Equal(t, 0, len(responseFinalScoreForUser.Body.Data), res)
	})

	t.Run("get final score for user with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)

		res, err := client.
			R().
			EnableTrace().
			SetQueryParam("user_played_quiz", userPlayedQuizId).
			Get("/api/v1/final_score/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseFinalScoreForUser utils.ResponseFinalScoreForUser
		err = json.Unmarshal(res.Body(), &responseFinalScoreForUser.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, 1, len(responseFinalScoreForUser.Body.Data), res)
		if len(responseFinalScoreForUser.Body.Data) == 1 {
			assert.Equal(t, "success", responseFinalScoreForUser.Body.Status, res)
			assert.Equal(t, 1, len(responseFinalScoreForUser.Body.Data), res)
			assert.Equal(t, 1, responseFinalScoreForUser.Body.Data[0].Rank, res)
			assert.Equal(t, "testcaseuser", responseFinalScoreForUser.Body.Data[0].UserName, res)
			assert.Equal(t, "testcaseuser", responseFinalScoreForUser.Body.Data[0].FirstName, res)
			assert.Equal(t, 0, responseFinalScoreForUser.Body.Data[0].Score, res)
			assert.Equal(t, -5, responseFinalScoreForUser.Body.Data[0].ResponseTime, res)
			assert.Equal(t, "Chase", responseFinalScoreForUser.Body.Data[0].ImageKey, res)
		}
	})
}
