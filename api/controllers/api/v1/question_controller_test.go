package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/stretchr/testify/assert"
)

var questionId string

// Define the expected result

func listQuestionsWithAnswerByQuizId(t *testing.T) {
	res, err := client.
		R().
		EnableTrace().
		Get(fmt.Sprintf("/api/v1/quizzes/%s/questions", quizId))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode(), res)

	var listSharedQuizzes utils.ResponseListQuestionByQuizId
	err = json.Unmarshal(res.Body(), &listSharedQuizzes.Body)
	assert.Nil(t, err, "Error in parsing response body")

	assert.Equal(t, 5, len(listSharedQuizzes.Body.Data.Data), res)
	questionId = listSharedQuizzes.Body.Data.Data[0].QuestionId
}

func TestListQuestionsWithAnswerByQuizId(t *testing.T) {
	t.Run("list questions with their answers with invalid input", func(t *testing.T) {
		invalidQuizId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"

		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/quizzes/%s/questions", invalidQuizId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("list questions with their answers with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		listQuestionsWithAnswerByQuizId(t)
	})
}

func TestGetQuestionById(t *testing.T) {
	t.Run("get question with their answers with invalid quiId", func(t *testing.T) {
		invalidQuizId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"

		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", invalidQuizId, questionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("get question with their answers with invalid questionID", func(t *testing.T) {
		invalidQuestionId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"
		createQuiz(t, quizTitle)

		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", quizId, invalidQuestionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode(), res)
	})

	t.Run("get questions with their answers with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		listQuestionsWithAnswerByQuizId(t)

		res, err := client.
			R().
			EnableTrace().
			Get(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", quizId, questionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var listSharedQuizzes utils.ResponseGetQuestionById
		err = json.Unmarshal(res.Body(), &listSharedQuizzes.Body)
		assert.Nil(t, err, "Error in parsing response body")

		expectedResult := structs.QuestionAnalytics{
			QuestionId:        questionId,
			CorrectAnswer:     "[2]",
			Question:          "Which city is known as the Eternal City?",
			Options:           map[string]string{"1": "Paris", "2": "Rome", "3": "Athens", "4": "Cairo"},
			QuestionsMedia:    "text",
			OptionsMedia:      "text",
			Resource:          "",
			Points:            1,
			QuestionTypeID:    1,
			QuestionType:      "single answer",
			DurationInSeconds: 60,
		}

		// Compare result
		assert.Equal(t, expectedResult.QuestionId, listSharedQuizzes.Body.Data.QuestionId, res)
		assert.Equal(t, expectedResult.Question, listSharedQuizzes.Body.Data.Question, res)
		assert.Equal(t, expectedResult.Options, listSharedQuizzes.Body.Data.Options, res)
		assert.Equal(t, expectedResult.QuestionsMedia, listSharedQuizzes.Body.Data.QuestionsMedia, res)
		assert.Equal(t, expectedResult.OptionsMedia, listSharedQuizzes.Body.Data.OptionsMedia, res)
		assert.Equal(t, expectedResult.Resource, listSharedQuizzes.Body.Data.Resource, res)
		assert.Equal(t, expectedResult.Points, listSharedQuizzes.Body.Data.Points, res)
		assert.Equal(t, expectedResult.QuestionTypeID, listSharedQuizzes.Body.Data.QuestionTypeID, res)
		assert.Equal(t, expectedResult.QuestionType, listSharedQuizzes.Body.Data.QuestionType, res)
		assert.Equal(t, expectedResult.DurationInSeconds, listSharedQuizzes.Body.Data.DurationInSeconds, res)
		assert.Equal(t, expectedResult.CorrectAnswer, listSharedQuizzes.Body.Data.CorrectAnswer, res)
	})
}

func TestUpdateQuestionById(t *testing.T) {
	t.Run("update question with invalid quiId", func(t *testing.T) {
		invalidQuizId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"

		res, err := client.
			R().
			EnableTrace().
			Put(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", invalidQuizId, questionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("update question with invalid input without body", func(t *testing.T) {
		createQuiz(t, quizTitle)
		listQuestionsWithAnswerByQuizId(t)

		res, err := client.
			R().
			EnableTrace().
			Put(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", quizId, questionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("update question with valid input", func(t *testing.T) {

		createQuiz(t, quizTitle)
		listQuestionsWithAnswerByQuizId(t)

		req := map[string]interface{}{
			"question":            "What is the capital of France?",
			"type":                1,
			"options":             map[string]string{"1": "Paris", "2": "London", "3": "Berlin", "4": "Madrid"},
			"answers":             []int{1},
			"points":              10,
			"duration_in_seconds": 30,
			"question_media":      "text",
			"options_media":       "text",
			"resource":            "",
		}

		res, err := client.
			R().
			EnableTrace().
			SetBody(req).
			Put(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", quizId, questionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}

func TestDeleteQuestionById(t *testing.T) {
	t.Run("delete question with invalid quiId", func(t *testing.T) {
		invalidQuizId := "fd40a699-a0ba-11ef-b98b-a4bb6d71ea0a"

		res, err := client.
			R().
			EnableTrace().
			Delete(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", invalidQuizId, questionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("delete question with invalid input with active quiz is present", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		listQuestionsWithAnswerByQuizId(t)

		res, err := client.
			R().
			EnableTrace().
			Delete(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", quizId, questionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("delete question with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		listQuestionsWithAnswerByQuizId(t)

		res, err := client.
			R().
			EnableTrace().
			Delete(fmt.Sprintf("/api/v1/quizzes/%s/questions/%s", quizId, questionId))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}
