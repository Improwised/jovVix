package utils

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/google/uuid"
	"github.com/jszwec/csvutil"
)

type Question struct {
	Question      string `csv:"Question Text"`
	Type          string `csv:"Question Type"`
	Points        string `csv:"Points,omitempty"`
	Option1       string `csv:"Option 1"`
	Option2       string `csv:"Option 2"`
	Option3       string `csv:"Option 3"`
	Option4       string `csv:"Option 4"`
	Option5       string `csv:"Option 5"`
	CorrectAnswer string `csv:"Correct Answer"`
	QuestionMedia string `csv:"Question Media"`
	OptionsMedia  string `csv:"Options Media"`
	Resource      string `csv:"Resource"`
}

func ValidateCSVFileFormat(fileName string) ([]Question, error) {
	var questions []Question

	// Open the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		return questions, err
	}
	defer file.Close()

	// Create a new CSV reader
	csvData, err := io.ReadAll(file)
	if err != nil {
		return questions, err
	}

	if err := csvutil.Unmarshal(csvData, &questions); err != nil {
		return questions, err
	}

	if len(questions) == 0 {
		return questions, fmt.Errorf(constants.ErrEmptyFile)
	}

	return questions, nil
}

func ExtractQuestionsFromCSV(questions []Question, questionTimeLimit string) ([]models.Question, error) {
	typeMapping := map[string]int{
		"single answer": 1,
		"survey":        2,
	}

	var validQuestions []models.Question
	for i, u := range questions {

		id, err := uuid.NewUUID()
		if err != nil {
			return validQuestions, err
		}

		options := make(map[string]string)
		if u.Option1 != "" {
			options["1"] = u.Option1
		}
		if u.Option2 != "" {
			options["2"] = u.Option2
		}
		if u.Option3 != "" {
			options["3"] = u.Option3
		}
		if u.Option4 != "" {
			options["4"] = u.Option4
		}
		if u.Option5 != "" {
			options["5"] = u.Option5
		}

		answers := []int{}
		for _, a := range strings.Split(u.CorrectAnswer, "|") {
			if a != "" {
				answerInt := 0
				fmt.Sscanf(a, "%d", &answerInt)
				answers = append(answers, answerInt)
			}
		}

		// Determine points
		points := 1
		if u.Points != "" {
			fmt.Sscanf(u.Points, "%d", &points)
		}

		duration := 30
		if questionTimeLimit != "" {
			if parsedDuration, err := strconv.Atoi(questionTimeLimit); err == nil {
				duration = parsedDuration
			}
		}

		validQuestions = append(validQuestions, models.Question{
			ID:                id,
			Question:          u.Question,
			Type:              typeMapping[u.Type],
			Options:           options,
			Answers:           answers,
			Points:            int16(points),
			DurationInSeconds: duration,
			OrderNumber:       i + 1,
			QuestionMedia:     u.QuestionMedia,
			OptionsMedia:      u.OptionsMedia,
			Resource:          sql.NullString{String: u.Resource, Valid: true},
		})
	}
	return validQuestions, nil
}
