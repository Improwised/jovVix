package utils

import (
	"sort"

	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
)

func ResponseToPdfData(quizReport []models.AnalyticsBoardAdmin) (map[int][]structs.UserAndQuestionData, []int) {
	orderToUserAndQuestionData := make(map[int][]structs.UserAndQuestionData)

	for _, data := range quizReport {
		user := structs.UserAndQuestionData{
			Options:        data.Options,
			UserName:       data.UserName,
			CorrectAnswer:  data.CorrectAnswer,
			SelectedAnswer: data.SelectedAnswer.String,
			Question:       data.Question,
			QuestionType:   data.QuestionType,
		}
		if val, ok := orderToUserAndQuestionData[data.OrderNo]; ok {
			orderToUserAndQuestionData[data.OrderNo] = append(val, user)
		} else {
			orderToUserAndQuestionData[data.OrderNo] = []structs.UserAndQuestionData{user}
		}
	}

	orderOfQuestion := make([]int, 0, len(orderToUserAndQuestionData))
	for o := range orderToUserAndQuestionData {
		orderOfQuestion = append(orderOfQuestion, o)
	}
	sort.Ints(orderOfQuestion)

	return orderToUserAndQuestionData, orderOfQuestion
}
