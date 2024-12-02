package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Improwised/quizz-app/api/logger"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type ResponseQuiz struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

var quizId string
var sessionId string
var quizTitle = "quiz for test cases"

func createQuiz(t *testing.T, quizTitle string) {
	res, err := client.
		R().
		EnableTrace().
		SetFile("attachment", "./controllers/api/v1/dummyCSVForTesting/demo.csv").
		SetFormData(map[string]string{
			"description": "This Quiz is create for test cases",
		}).
		Post(fmt.Sprintf("/api/v1/quizzes/%s/upload", quizTitle))

	assert.Nil(t, err)
	assert.Equal(t, http.StatusAccepted, res.StatusCode(), res)

	var quizResponse ResponseQuiz
	err = json.Unmarshal(res.Body(), &quizResponse)
	assert.Nil(t, err, "Error in parsing response body")

	quizId = quizResponse.Data
}

func generateDemoSession(t *testing.T) {

	res, err := client.
		R().
		EnableTrace().
		Post(fmt.Sprintf("/api/v1/quizzes/%s/demo_session", quizId))

	assert.Nil(t, err)
	assert.Equal(t, http.StatusAccepted, res.StatusCode(), res)

	var sessionResponse ResponseQuiz
	err = json.Unmarshal(res.Body(), &sessionResponse)
	assert.Nil(t, err, "Error in parsing response body")

	sessionId = sessionResponse.Data
}

func terminateQuiz(t *testing.T) {
	logger, err := logger.NewRootLogger(true, true)
	assert.Nil(t, err)

	parsedSessionId, err := uuid.Parse(sessionId)
	assert.Nil(t, err)

	activeQuizModel := models.InitActiveQuizModel(db, logger)

	err = activeQuizModel.Deactivate(parsedSessionId)
	assert.Nil(t, err)
}

func TestGenerateDemoSession(t *testing.T) {
	invalidQuizId := "00000000-5a14-11ef-866f-a4bb6d71ea0q"
	t.Run("generate demo session with invalid quizId", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Post(fmt.Sprintf("/api/v1/quizzes/%s/demo_session", invalidQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("generate demo session with valid quizId", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
	})
}

func TestCreateQuizByCsv(t *testing.T) {
	t.Run("create quiz wit invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			SetFile("attachment", "./controllers/api/v1/dummyCSVForTesting/EmptyFile.csv").
			SetFormData(map[string]string{
				"description": "This Quiz is create for test cases",
			}).
			Post(fmt.Sprintf("/api/v1/quizzes/%s/upload", quizTitle))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("create quiz wit valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
	})
}

func TestGetAdminUploadedQuizzes(t *testing.T) {

	t.Run("get uploaded quiz wit valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)

		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/quizzes")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var quizzesResponse utils.ResponseAdminUploadedQuiz
		err = json.Unmarshal(res.Body(), &quizzesResponse.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.GreaterOrEqual(t, len(quizzesResponse.Body.Data), 1, "Expected at least one quiz for the admin")
	})
}

func TestDeleteQuizById(t *testing.T) {

	t.Run("deletle quiz with invalid input (active quiz is present)", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)

		res, err := client.
			R().
			EnableTrace().
			Delete(fmt.Sprintf("/api/v1/quizzes/%s", quizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("deletle quiz with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)

		res, err := client.
			R().
			EnableTrace().
			Delete(fmt.Sprintf("/api/v1/quizzes/%s", quizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}

func TestGetQuizAnalysis(t *testing.T) {

	t.Run("get quiz analysis for admin with unauthorized user", func(t *testing.T) {
		res, err := userclient.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/admin/reports/%s/analysis", invaliduserPlayedQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("get quiz analysis for admin with invalid input (active_quiz_id is invalid)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/admin/reports/%s/analysis", invaliduserPlayedQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseGetQuizAnalysis utils.ResponseGetQuizAnalysis
		err = json.Unmarshal(res.Body(), &responseGetQuizAnalysis.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "success", responseGetQuizAnalysis.Body.Status, res)
		assert.Equal(t, 0, len(responseGetQuizAnalysis.Body.Data), res)
	})

	t.Run("get quiz analysis for admin with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)

		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/admin/reports/%s/analysis", sessionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseGetQuizAnalysis utils.ResponseGetQuizAnalysis
		err = json.Unmarshal(res.Body(), &responseGetQuizAnalysis.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, 5, len(responseGetQuizAnalysis.Body.Data), res)
		if len(responseGetQuizAnalysis.Body.Data) >= 1 {
			assert.Equal(t, "success", responseGetQuizAnalysis.Body.Status, res)
			assert.Equal(t, 5, len(responseGetQuizAnalysis.Body.Data), res)
			assert.Equal(t, "Which city is known as the Eternal City?", responseGetQuizAnalysis.Body.Data[0].Question, res)
			assert.Equal(t, 1, responseGetQuizAnalysis.Body.Data[0].Type, res)
			assert.Equal(t, map[string]string{"1": "Paris", "2": "Rome", "3": "Athens", "4": "Cairo"}, responseGetQuizAnalysis.Body.Data[0].Options, res)
			assert.Equal(t, "text", responseGetQuizAnalysis.Body.Data[0].QuestionsMedia, res)
			assert.Equal(t, "text", responseGetQuizAnalysis.Body.Data[0].OptionsMedia, res)
			assert.Equal(t, "", responseGetQuizAnalysis.Body.Data[0].Resource, res)
			assert.Equal(t, []int{2}, responseGetQuizAnalysis.Body.Data[0].CorrectAnswers, res)
			assert.Equal(t, map[string]interface{}{"testcaseuser": nil}, responseGetQuizAnalysis.Body.Data[0].SelectedAnswers, res)
			assert.Equal(t, 60, responseGetQuizAnalysis.Body.Data[0].DurationInSeconds, res)
			assert.Equal(t, float32(-1), responseGetQuizAnalysis.Body.Data[0].AvgResponseTime, res)
		}
	})
}

func TestListQuizzesAnalysis(t *testing.T) {

	t.Run("list quizzes analysis for admin with unauthorized user", func(t *testing.T) {
		res, err := userclient.
			R().
			EnableTrace().
			Get("/api/v1/admin/reports/list")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("list quizzes analysis for admin with invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/admin/reports/list")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseListQuizzesAnalysis utils.RsponseListQuizzesAnalysis
		err = json.Unmarshal(res.Body(), &responseListQuizzesAnalysis.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "success", responseListQuizzesAnalysis.Body.Status, res)
		assert.Equal(t, 0, len(responseListQuizzesAnalysis.Body.Data.Data), res)
		assert.Equal(t, int64(0), responseListQuizzesAnalysis.Body.Data.Count, res)
	})

	t.Run("list quizzes analysis for admin with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)
		terminateQuiz(t)

		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/admin/reports/list")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseListQuizzesAnalysis utils.RsponseListQuizzesAnalysis
		err = json.Unmarshal(res.Body(), &responseListQuizzesAnalysis.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, 1, len(responseListQuizzesAnalysis.Body.Data.Data), res)
		if len(responseListQuizzesAnalysis.Body.Data.Data) >= 1 {
			assert.Equal(t, "success", responseListQuizzesAnalysis.Body.Status, res)
			assert.Equal(t, 1, len(responseListQuizzesAnalysis.Body.Data.Data), res)
			assert.Equal(t, "quiz%20for%20test%20cases", responseListQuizzesAnalysis.Body.Data.Data[0].Title, res)
			assert.Equal(t, "This Quiz is create for test cases", responseListQuizzesAnalysis.Body.Data.Data[0].Description.String, res)
			assert.Equal(t, 5, responseListQuizzesAnalysis.Body.Data.Data[0].Questions, res)
			assert.Equal(t, 1, responseListQuizzesAnalysis.Body.Data.Data[0].Participants, res)
			assert.Equal(t, 0, responseListQuizzesAnalysis.Body.Data.Data[0].CorrectAnswers, res)
		}
	})
}
