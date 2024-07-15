package quizUtilsHelper

import (
	"errors"

	"github.com/Improwised/quizz-app/api/constants"
)

/*
1 - Single Answer
2 - Survey
*/

// add other types to first constants and then here
var questionTypeIDs = map[string]int{
	constants.SingleAnswerString: constants.SingleAnswer,
	constants.SurveyString:       constants.Survey,
}

// function to check if passed type exist as a type or not
func CheckQuestionType(questionType string) (int, error) {
	if id, ok := questionTypeIDs[questionType]; ok {
		return id, nil
	}
	return 0, errors.New(constants.ErrQuestionType)
}

// function to get the string when number is provided
func GetQuestionType(typeNumber int) (string, error) {
	for questionType, id := range questionTypeIDs {
		if id == typeNumber {
			return questionType, nil
		}
	}
	return "", errors.New(constants.ErrQuestionId)
}
