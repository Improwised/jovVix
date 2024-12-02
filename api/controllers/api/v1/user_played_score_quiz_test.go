package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Improwised/quizz-app/api/logger"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/stretchr/testify/assert"
)

var invalidInvitationCode = 62674
var userPlayedQuizId string
var invaliduserPlayedQuizId = "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"

func playedQuizValidation(t *testing.T) {
	logger, err := logger.NewRootLogger(true, true)
	assert.Nil(t, err)
	activeQuizModel := models.InitActiveQuizModel(db, logger)

	createQuiz(t, quizTitle)
	generateDemoSession(t)

	session, err := activeQuizModel.GetOrActivateSession(sessionId, userId)
	assert.Nil(t, err)

	res, err := userclient.
		R().
		EnableTrace().
		Post(fmt.Sprintf("/api/v1/user_played_quizes/%d", session.InvitationCode.Int32))

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode(), res)

	var responsePlayedQuizValidation utils.ResponsePlayedQuizValidation
	err = json.Unmarshal(res.Body(), &responsePlayedQuizValidation.Body)
	assert.Nil(t, err, "Error in parsing response body")

	assert.Equal(t, sessionId, responsePlayedQuizValidation.Body.Data.SessionId, res)
	userPlayedQuizId = responsePlayedQuizValidation.Body.Data.UserPlayedQuizId
}

func TestPlayedQuizValidation(t *testing.T) {

	t.Run("check played quiz validation with invalid input", func(t *testing.T) {
		res, err := userclient.
			R().
			EnableTrace().
			Post(fmt.Sprintf("/api/v1/user_played_quizes/%d", invalidInvitationCode))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("check played quiz validation with invalid input (host cannot be a player in their own quiz)", func(t *testing.T) {
		logger, err := logger.NewRootLogger(true, true)
		assert.Nil(t, err)
		activeQuizModel := models.InitActiveQuizModel(db, logger)

		createQuiz(t, quizTitle)
		generateDemoSession(t)

		session, err := activeQuizModel.GetOrActivateSession(sessionId, userId)
		assert.Nil(t, err)

		res, err := client.
			R().
			EnableTrace().
			Post(fmt.Sprintf("/api/v1/user_played_quizes/%d", session.InvitationCode.Int32))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode(), res)
	})

	t.Run("check played quiz validation with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)
	})
}

func TestListUserPlayedQuizes(t *testing.T) {

	t.Run("list user played quizzes with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)

		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/user_played_quizes")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		fmt.Println(res)
		var responseListUserPlayedQuizes utils.ResponseListUserPlayedQuizes
		err = json.Unmarshal(res.Body(), &responseListUserPlayedQuizes.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, int64(0), responseListUserPlayedQuizes.Body.Data.Count, res)
	})
}

func TestListUserPlayedQuizesWithQuestionById(t *testing.T) {

	t.Run("list user played quizzes with invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/user_played_quizes/%s", invaliduserPlayedQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseuserPlayedQuizesWithQuestion utils.ResponseListUserPlayedQuizesWithQuestionById
		err = json.Unmarshal(res.Body(), &responseuserPlayedQuizesWithQuestion.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "success", responseuserPlayedQuizesWithQuestion.Body.Status, res)
		assert.Equal(t, 0, len(responseuserPlayedQuizesWithQuestion.Body.Data), res)
	})

	t.Run("list user played quizzes with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)

		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/user_played_quizes/%s", userPlayedQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseuserPlayedQuizesWithQuestion utils.ResponseListUserPlayedQuizesWithQuestionById
		err = json.Unmarshal(res.Body(), &responseuserPlayedQuizesWithQuestion.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "success", responseuserPlayedQuizesWithQuestion.Body.Status, res)
		assert.Equal(t, 5, len(responseuserPlayedQuizesWithQuestion.Body.Data), res)
	})
}
