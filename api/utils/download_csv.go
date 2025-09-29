package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
)

type InputParticipants struct {
	QuizAnalysis         []models.QuizAnalysis
	FinalScoreBoardAdmin []models.FinalScoreBoardAdmin
}

type UserStatics struct {
	Accuracy       int
	CorrectAnswers int
	SurveyAnswers  int
	Attempted      int
}

func GetQuestionsData(analysis []models.QuizAnalysis, questionType string) [][]string {
	questions := GetQuestionsByType(analysis, questionType)

	maxOptions := MaximumCountOfOptions(questions)

	// header
	header := []string{
		"Question Text",
		"Question Type",
	}
	for i := 1; i <= maxOptions; i++ {
		header = append(header, fmt.Sprintf("Option %d", i))
		header = append(header, fmt.Sprintf("Option %d count", i))
	}
	header = append(header, "Correct Answer", "Correct Answer Percentage")

	formatedQuestionData := [][]string{header}

	// rows
	for _, question := range questions {
		options := GetOptions(question)
		correctAnsField := ""

		row := []string{question.Question}

		switch questionType {
		case constants.AllString:
			switch question.Type {
			case 1:
				row = append(row, constants.SingleAnsField)
			case 2:
				row = append(row, constants.SurveyString)
			}
		case constants.SingleString:
			row = append(row, constants.SingleAnsField)
		case constants.SurveyString:
			row = append(row, constants.SurveyString)
		}

		for _, option := range options {
			row = append(row, option.Option, fmt.Sprintf("%d", option.SelectedCount))
		}

		missing := maxOptions - len(options)
		for i := 0; i < missing; i++ {
			row = append(row, "", "")
		}

		for _, correctAns := range question.CorrectAnswers {
			correctAnsField += fmt.Sprintf("%d|", correctAns)
		}
		correctAnsField = strings.TrimRight(correctAnsField, "|")

		row = append(row, correctAnsField, fmt.Sprintf("%.2f", CountPercentageOfGivenAns(question)))

		formatedQuestionData = append(formatedQuestionData, row)
	}

	return formatedQuestionData
}

func GetParticipantsData(analysis InputParticipants) [][]string {
	userStats := GetUserStatics(analysis.QuizAnalysis)

	// CSV header
	records := [][]string{
		{"Rank", "UserName", "Score", "Accuracy", "CorrectAnswers", "SurveyAnswers", "Attempted"},
	}

	for _, u := range analysis.FinalScoreBoardAdmin {
		accuracy, correct, survey, attempted := 0, 0, 0, 0

		if stat, ok := userStats[u.UserName]; ok {
			accuracy = stat.Accuracy
			correct = stat.CorrectAnswers
			survey = stat.SurveyAnswers
			attempted = stat.Attempted
		}

		record := []string{
			strconv.Itoa(u.Rank),
			u.UserName,
			strconv.Itoa(u.Score),
			strconv.Itoa(accuracy),
			strconv.Itoa(correct),
			strconv.Itoa(survey),
			strconv.Itoa(attempted),
		}

		records = append(records, record)
	}

	return records
}

func GetQuestionsByType(analysis []models.QuizAnalysis, questionType string) []models.QuizAnalysis {
	var questions []models.QuizAnalysis

	if questionType == constants.SingleString {
		for _, question := range analysis {
			if question.Type == constants.SingleAnswer {
				questions = append(questions, question)
			}
		}
		return questions
	}
	if questionType == constants.SurveyString {
		for _, question := range analysis {
			if question.Type == constants.Survey {
				questions = append(questions, question)
			}
		}
		return questions
	}
	return analysis
}

func GetOptions(question models.QuizAnalysis) []structs.OptionsWithSelectedCount {
	questionOptions := make(map[string]structs.OptionsWithSelectedCount)

	for key, val := range question.Options {
		questionOptions[key] = structs.OptionsWithSelectedCount{
			Option:        val,
			SelectedCount: 0,
		}
	}

	for _, raw := range question.SelectedAnswers {
		if selected, ok := raw.([]interface{}); ok {
			for _, s := range selected {
				if optFloat, ok := s.(float64); ok {
					optKey := fmt.Sprintf("%d", int(optFloat))
					if option, exists := questionOptions[optKey]; exists {
						option.SelectedCount++
						questionOptions[optKey] = option
					}
				}
			}
		}
	}

	orderedOptions := []structs.OptionsWithSelectedCount{}
	for i := 1; i <= len(question.Options); i++ {
		key := fmt.Sprintf("%d", i)
		if opt, ok := questionOptions[key]; ok {
			orderedOptions = append(orderedOptions, opt)
		}
	}

	return orderedOptions
}

func CountPercentageOfGivenAns(question models.QuizAnalysis) float64 {
	totalUser := len(question.SelectedAnswers)
	correctCount := 0

	for _, correct := range question.CorrectAnswers {
		for _, selectedAns := range question.SelectedAnswers {
			if arr, ok := selectedAns.([]interface{}); ok {
				for _, v := range arr {
					if option, ok := v.(float64); ok {
						if int(option) == correct {
							correctCount++
						}
					}
				}
			}

		}
	}

	totalPercentage := (float64(correctCount) / float64(totalUser)) * 100
	return totalPercentage
}

func MaximumCountOfOptions(analysis []models.QuizAnalysis) int {
	maxNumberOfOption := 0
	for _, que := range analysis {
		lenOfOptions := len(que.Options)
		if lenOfOptions > maxNumberOfOption {
			maxNumberOfOption = lenOfOptions
		}
	}
	return maxNumberOfOption
}

func GetUserStatics(analysis []models.QuizAnalysis) map[string]UserStatics {
	totalQuestion := 0
	userStats := make(map[string]UserStatics)

	for _, q := range analysis {

		for user, selectedRaw := range q.SelectedAnswers {
			stat := userStats[user]

			var selected []int
			if arr, ok := selectedRaw.([]interface{}); ok {
				for _, v := range arr {
					if num, ok := v.(float64); ok {
						selected = append(selected, int(num))
					}
				}
			}

			if len(selected) > 0 {
				stat.Attempted++
			}

			if q.Type == 2 {
				stat.SurveyAnswers++
			} else {
				if isCorrectAnswer(selected, q.CorrectAnswers) {
					stat.CorrectAnswers++
				}
			}

			userStats[user] = stat
		}
		totalQuestion++
	}

	for user, stat := range userStats {
		if totalQuestion > 0 {
			stat.Accuracy = int((float64(stat.CorrectAnswers+stat.SurveyAnswers) / float64(totalQuestion)) * 100)
		}
		userStats[user] = stat
	}

	return userStats
}

func isCorrectAnswer(selected, correct []int) bool {
	if len(selected) == 0 || len(selected) != len(correct) {
		return false
	}
	sort.Ints(selected)
	sort.Ints(correct)
	for i := range selected {
		if selected[i] != correct[i] {
			return false
		}
	}
	return true
}
