package v1_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertImage(t *testing.T) {
	t.Run("insert images with invalid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			SetFile("image-attachment", "./controllers/api/v1/dummyCSVForTesting/avatar.png").
			Post("/api/v1/images")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("insert images with invalid input (not give filename as questionId)", func(t *testing.T) {
		createQuiz(t, quizTitle)

		res, err := client.
			R().
			EnableTrace().
			SetQueryParam("quiz_id", quizId).
			SetFile("image-attachment", "./controllers/api/v1/dummyCSVForTesting/avatar.png").
			Post("/api/v1/images")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode(), res)
	})

	t.Run("insert images with invalid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		listQuestionsWithAnswerByQuizId(t)

		tempFilePath := "./controllers/api/v1/dummyCSVForTesting/" + questionId
		err := os.Rename("./controllers/api/v1/dummyCSVForTesting/avatar.png", tempFilePath)
		if err != nil {
			t.Fatalf("failed to rename file: %v", err)
		}
		defer os.Rename(tempFilePath, "./controllers/api/v1/dummyCSVForTesting/avatar.png")

		res, err := client.
			R().
			EnableTrace().
			SetQueryParam("quiz_id", quizId).
			SetFile("image-attachment", tempFilePath).
			Post("/api/v1/images")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}
