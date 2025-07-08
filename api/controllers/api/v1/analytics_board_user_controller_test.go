package v1_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAnalyticsForUser(t *testing.T) {

	t.Run("get analytics score for user with invalid input (user_played_quiz query params is not given)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/analytics_board/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("get analytics score for user with invalid input (user_played_quiz query params is invalid)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			SetQueryParam(constants.UserPlayedQuiz, invaliduserPlayedQuizId).
			Get("/api/v1/analytics_board/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseAnalyticsBoardForUser utils.ResponseAnalyticsBoardForUser
		err = json.Unmarshal(res.Body(), &responseAnalyticsBoardForUser.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "success", responseAnalyticsBoardForUser.Body.Status, res)
		assert.Equal(t, 0, len(responseAnalyticsBoardForUser.Body.Data), res)
	})

	t.Run("get analytics score for user with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)

		res, err := client.
			R().
			EnableTrace().
			SetQueryParam(constants.UserPlayedQuiz, userPlayedQuizId).
			Get("/api/v1/analytics_board/user")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseAnalyticsBoardForUser utils.ResponseAnalyticsBoardForUser
		err = json.Unmarshal(res.Body(), &responseAnalyticsBoardForUser.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, 5, len(responseAnalyticsBoardForUser.Body.Data), res)
		if len(responseAnalyticsBoardForUser.Body.Data) >= 1 {
			assert.Equal(t, "success", responseAnalyticsBoardForUser.Body.Status, res)
			assert.Equal(t, 5, len(responseAnalyticsBoardForUser.Body.Data), res)
			assert.Equal(t, "testcaseuser", responseAnalyticsBoardForUser.Body.Data[0].UserName, res)
			assert.Equal(t, "testcaseuser", responseAnalyticsBoardForUser.Body.Data[0].FirstName, res)
			assert.Equal(t, "[2]", responseAnalyticsBoardForUser.Body.Data[0].CorrectAnswer, res)
			assert.Equal(t, 0, responseAnalyticsBoardForUser.Body.Data[0].CalculatedScore, res)
			assert.Equal(t, false, responseAnalyticsBoardForUser.Body.Data[0].IsAttend, res)
			assert.Equal(t, -1, responseAnalyticsBoardForUser.Body.Data[0].ResponseTime, res)
			assert.Equal(t, 0, responseAnalyticsBoardForUser.Body.Data[0].CalculatedPoints, res)
			assert.Equal(t, "Which city is known as the Eternal City?", responseAnalyticsBoardForUser.Body.Data[0].Question, res)
			assert.Equal(t, map[string]string{"1": "Paris", "2": "Rome", "3": "Athens", "4": "Cairo"}, responseAnalyticsBoardForUser.Body.Data[0].Options, res)
			assert.Equal(t, "text", responseAnalyticsBoardForUser.Body.Data[0].QuestionsMedia, res)
			assert.Equal(t, "text", responseAnalyticsBoardForUser.Body.Data[0].OptionsMedia, res)
			assert.Equal(t, "", responseAnalyticsBoardForUser.Body.Data[0].Resource, res)
			assert.Equal(t, 1, responseAnalyticsBoardForUser.Body.Data[0].Points, res)
			assert.Equal(t, 1, responseAnalyticsBoardForUser.Body.Data[0].QuestionTypeID, res)
			assert.Equal(t, "single answer", responseAnalyticsBoardForUser.Body.Data[0].QuestionType, res)
			assert.Equal(t, 1, responseAnalyticsBoardForUser.Body.Data[0].OrderNo, res)
		}
	})
}
