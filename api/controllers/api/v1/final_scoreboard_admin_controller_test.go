package v1_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetScoreForAdmin(t *testing.T) {

	t.Run("get final score for admin with unauthorized user", func(t *testing.T) {
		res, err := userclient.
			R().
			EnableTrace().
			Get("/api/v1/final_score/admin")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("get final score for admin with invalid input (active_quiz_id query params is not given)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/final_score/admin")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("get final score for admin with invalid input (active_quiz_id query params is invalid)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			SetQueryParam(constants.ActiveQuizId, invaliduserPlayedQuizId).
			Get("/api/v1/final_score/admin")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseFinalScoreForAdmin utils.ResponseFinalScoreForAdmin
		err = json.Unmarshal(res.Body(), &responseFinalScoreForAdmin.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "success", responseFinalScoreForAdmin.Body.Status, res)
		assert.Equal(t, 0, len(responseFinalScoreForAdmin.Body.Data), res)
	})

	t.Run("get final score for admin with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)

		res, err := client.
			R().
			EnableTrace().
			SetQueryParam(constants.ActiveQuizId, sessionId).
			Get("/api/v1/final_score/admin")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseFinalScoreForAdmin utils.ResponseFinalScoreForAdmin
		err = json.Unmarshal(res.Body(), &responseFinalScoreForAdmin.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, 1, len(responseFinalScoreForAdmin.Body.Data), res)
		if len(responseFinalScoreForAdmin.Body.Data) == 1 {
			assert.Equal(t, "success", responseFinalScoreForAdmin.Body.Status, res)
			assert.Equal(t, 1, len(responseFinalScoreForAdmin.Body.Data), res)
			assert.Equal(t, 1, responseFinalScoreForAdmin.Body.Data[0].Rank, res)
			assert.Equal(t, "testcaseuser", responseFinalScoreForAdmin.Body.Data[0].UserName, res)
			assert.Equal(t, "testcaseuser", responseFinalScoreForAdmin.Body.Data[0].FirstName, res)
			assert.Equal(t, 0, responseFinalScoreForAdmin.Body.Data[0].Score, res)
			assert.Equal(t, -5, responseFinalScoreForAdmin.Body.Data[0].ResponseTime, res)
			assert.Equal(t, "Chase", responseFinalScoreForAdmin.Body.Data[0].ImageKey, res)
		}
	})
}
