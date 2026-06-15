package utils

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
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

// normalizeMedia trims and lowercases a media value and reports whether it is
// allowed. An empty value defaults to "text" to match manual question creation.
func normalizeMedia(media string) (string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(media))
	if normalized == "" {
		return constants.MediaText, true
	}
	switch normalized {
	case constants.MediaText, constants.MediaImage, constants.MediaCode:
		return normalized, true
	default:
		return normalized, false
	}
}

func ExtractQuestionsFromCSV(questions []Question, questionTimeLimit string) ([]models.Question, error) {
	// Duration comes solely from configuration; reject if it is missing or invalid.
	duration, err := strconv.Atoi(strings.TrimSpace(questionTimeLimit))
	if err != nil || duration <= 0 {
		return nil, fmt.Errorf(constants.ErrInvalidQuestionTimeLimit)
	}

	var validQuestions []models.Question
	var rowErrors []string

	for i, u := range questions {
		// Row number as seen by the user in a spreadsheet (header is row 1).
		rowNo := i + 2

		var rowIssues []string

		// Question text must be present.
		if strings.TrimSpace(u.Question) == "" {
			rowIssues = append(rowIssues, constants.ErrEmptyQuestionText)
		}

		// Question type must be a known type (single answer / survey).
		questionType, typeErr := quizUtilsHelper.CheckQuestionType(strings.TrimSpace(u.Type))
		if typeErr != nil {
			rowIssues = append(rowIssues, fmt.Sprintf("%s (got %q, allowed: %s, %s)", constants.ErrQuestionType, u.Type, constants.SingleAnswerString, constants.SurveyString))
		}

		// Collect non-empty options, preserving their option number.
		options := make(map[string]string)
		for idx, opt := range []string{u.Option1, u.Option2, u.Option3, u.Option4, u.Option5} {
			if strings.TrimSpace(opt) != "" {
				options[strconv.Itoa(idx+1)] = opt
			}
		}
		if len(options) < 2 {
			rowIssues = append(rowIssues, constants.ErrInsufficientOptions)
		}

		// Correct answer(s): must be present, numeric, and reference an existing option.
		answers := []int{}
		correctRaw := strings.TrimSpace(u.CorrectAnswer)
		if correctRaw == "" {
			rowIssues = append(rowIssues, constants.ErrEmptyCorrectAnswer)
		} else {
			for _, a := range strings.Split(correctRaw, "|") {
				a = strings.TrimSpace(a)
				if a == "" {
					continue
				}
				answerInt, convErr := strconv.Atoi(a)
				if convErr != nil {
					rowIssues = append(rowIssues, fmt.Sprintf("%s (got %q)", constants.ErrInvalidCorrectAnswer, a))
					continue
				}
				if _, ok := options[strconv.Itoa(answerInt)]; !ok {
					rowIssues = append(rowIssues, fmt.Sprintf("%s (option %d does not exist)", constants.ErrInvalidCorrectAnswer, answerInt))
					continue
				}
				answers = append(answers, answerInt)
			}
		}

		// Type-specific answer count rules (only meaningful when the type is valid).
		if typeErr == nil {
			switch questionType {
			case constants.SingleAnswer:
				if len(answers) != 1 {
					rowIssues = append(rowIssues, constants.ErrSingleAnswerLength)
				}
			case constants.Survey:
				if len(answers) < 1 {
					rowIssues = append(rowIssues, constants.ErrSurveyAnswerLength)
				}
			}
		}

		// Media types: optional (default text), but must be text, image, or code.
		questionMedia, questionMediaOK := normalizeMedia(u.QuestionMedia)
		if !questionMediaOK {
			rowIssues = append(rowIssues, fmt.Sprintf("%s (got %q)", constants.ErrInvalidQuestionMedia, u.QuestionMedia))
		}
		optionsMedia, optionsMediaOK := normalizeMedia(u.OptionsMedia)
		if !optionsMediaOK {
			rowIssues = append(rowIssues, fmt.Sprintf("%s (got %q)", constants.ErrInvalidOptionsMedia, u.OptionsMedia))
		}

		// Points: optional, but if present must be a positive integer.
		points := 1
		if strings.TrimSpace(u.Points) != "" {
			parsedPoints, convErr := strconv.Atoi(strings.TrimSpace(u.Points))
			if convErr != nil || parsedPoints <= 0 {
				rowIssues = append(rowIssues, fmt.Sprintf("%s (got %q)", constants.ErrInvalidPoints, u.Points))
			} else {
				points = parsedPoints
			}
		}

		if len(rowIssues) > 0 {
			rowErrors = append(rowErrors, fmt.Sprintf("row %d: %s", rowNo, strings.Join(rowIssues, "; ")))
			continue
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return validQuestions, err
		}

		validQuestions = append(validQuestions, models.Question{
			ID:                id,
			Question:          u.Question,
			Type:              questionType,
			Options:           options,
			Answers:           answers,
			Points:            int16(points),
			DurationInSeconds: duration,
			OrderNumber:       i + 1,
			QuestionMedia:     questionMedia,
			OptionsMedia:      optionsMedia,
			Resource:          sql.NullString{String: u.Resource, Valid: true},
		})
	}

	if len(rowErrors) > 0 {
		return nil, fmt.Errorf("%s %s", constants.ErrInvalidCSVRows, strings.Join(rowErrors, " | "))
	}

	return validQuestions, nil
}
