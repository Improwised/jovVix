package quizUtilsHelper

import (
	"testing"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/stretchr/testify/assert"
)

func TestCheckQuestionType(t *testing.T) {
	// test for static string
	questionType := "single answer"

	questionID, err := CheckQuestionType(questionType)

	assert.NoError(t, err)
	assert.Equal(t, constants.SingleAnswer, questionID, "Expected result to be %d, but got %d", constants.SingleAnswer, questionID)

	// test for string from constants
	questionType = constants.SurveyString

	questionID, err = CheckQuestionType(questionType)

	assert.NoError(t, err)
	assert.Equal(t, constants.Survey, questionID, "Expected result to be %d, but got %d", constants.Survey, questionID)

	// test for a type which is not there in the map
	questionType = "multi answer"

	questionID, err = CheckQuestionType(questionType)

	assert.Error(t, err)
	assert.Equal(t, 0, questionID, "Expected result to be %d, but got %d", 0, questionID)
}

func TestGetTypeNumber(t *testing.T) {
	typeNumber := constants.SingleAnswer

	questionType, err := GetQuestionType(typeNumber)

	assert.NoError(t, err)
	assert.Equal(t, constants.SingleAnswerString, questionType, "Expected result to be %s, but got %s", constants.SingleAnswerString, questionType)

	typeNumber = constants.Survey

	questionType, err = GetQuestionType(typeNumber)

	assert.NoError(t, err)
	assert.Equal(t, constants.SurveyString, questionType, "Expected result to be %s, but got %s", constants.SurveyString, questionType)

	typeNumber = 145

	questionType, err = GetQuestionType(typeNumber)

	assert.Error(t, err)
	assert.Equal(t, "", questionType, "Expected result to be %s, but got %s", "", questionType)

	typeNumber = 0

	questionType, err = GetQuestionType(typeNumber)

	assert.Error(t, err)
	assert.Equal(t, "", questionType, "Expected result to be %s, but got %s", "", questionType)
}
