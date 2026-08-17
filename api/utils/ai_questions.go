package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
)

func normalizeAIQuestionType(questionType string) (int, bool) {
	switch strings.ToLower(strings.TrimSpace(questionType)) {
	case "":
		return constants.SingleAnswer, true
	case constants.AIQuestionTypeSurvey:
		return constants.Survey, true
	case constants.AIQuestionTypeSingle, constants.SingleAnswerString, "mcq", "multiple choice":
		return constants.SingleAnswer, true
	default:
		return constants.SingleAnswer, false
	}
}

func aiQuestionTypeWord(questionType int) string {
	if questionType == constants.Survey {
		return constants.AIQuestionTypeSurvey
	}
	return constants.AIQuestionTypeSingle
}

type AIGeneration struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Questions   []structs.AIQuestion `json:"questions"`
}

type AIQuestionOptions struct {
	Duration     int
	Points       int16
	MaxQuestions int
}

type AINormalizeResult struct {
	Questions []models.Question
	Accepted  []structs.AIQuestion
	Issues    []string
}

func ExtractAIGeneration(content string) (AIGeneration, error) {
	payload := stripCodeFence(content)

	if object := sliceBetween(payload, '{', '}'); object != "" {
		var generation AIGeneration
		if err := json.Unmarshal([]byte(object), &generation); err == nil && len(generation.Questions) > 0 {
			return generation, nil
		}
	}

	if array := sliceBetween(payload, '[', ']'); array != "" {
		var questions []structs.AIQuestion
		if err := json.Unmarshal([]byte(array), &questions); err == nil && len(questions) > 0 {
			return AIGeneration{Questions: questions}, nil
		}
	}

	return AIGeneration{}, errors.New(constants.ErrAIInvalidResponse)
}

func SanitizeAITitle(title, topic string) string {
	collapsed := strings.Join(strings.Fields(title), " ")
	if collapsed == "" {
		collapsed = strings.TrimSpace(strings.Join(strings.Fields(topic), " ") + " quiz")
	}
	return TruncateRunes(collapsed, constants.AIMaxTitleLength)
}

func SanitizeAIDescription(description string) string {
	return TruncateRunes(strings.Join(strings.Fields(description), " "), constants.AIMaxDescriptionLength)
}

// TruncateRunes cuts a string to a rune limit, not a byte limit.
func TruncateRunes(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return strings.TrimSpace(string(runes[:limit]))
}

func NormalizeAIQuestions(raw []structs.AIQuestion, opts AIQuestionOptions) (AINormalizeResult, error) {
	var result AINormalizeResult

	if opts.Duration <= 0 {
		return result, errors.New(constants.ErrInvalidQuestionTimeLimit)
	}

	maxQuestions := opts.MaxQuestions
	if maxQuestions <= 0 {
		maxQuestions = constants.AIDefaultMaxQuestions
	}

	points := opts.Points
	if points <= constants.MinimumPoints {
		points = constants.AIDefaultQuestionPoints
	}
	if points > constants.MaximumPoints {
		points = constants.MaximumPoints
	}

	seen := make(map[string]bool)
	truncated := false

	for i, item := range raw {
		if len(result.Questions) >= maxQuestions {
			truncated = true
			break
		}

		questionNo := i + 1
		var issues []string

		question := strings.TrimSpace(item.Question)
		if question == "" {
			issues = append(issues, constants.ErrEmptyQuestionText)
		} else if utf8.RuneCountInString(question) > constants.AIMaxQuestionLength {
			issues = append(issues, constants.ErrAIQuestionTooLong)
		}

		questionMedia, questionMediaOK := normalizeMedia(item.QuestionMedia)
		if !questionMediaOK || questionMedia == constants.MediaImage {
			issues = append(issues, fmt.Sprintf("%s (got %q)", constants.ErrAIUnsupportedQuestionMedia, item.QuestionMedia))
			questionMedia = constants.MediaText
		}

		optionsMedia, optionsMediaOK := normalizeMedia(item.OptionsMedia)
		if !optionsMediaOK || optionsMedia == constants.MediaImage {
			issues = append(issues, fmt.Sprintf("%s (got %q)", constants.ErrAIUnsupportedOptionsMedia, item.OptionsMedia))
			optionsMedia = constants.MediaText
		}

		resource := strings.TrimSpace(stripCodeFence(item.Resource))
		if questionMedia == constants.MediaCode {
			if resource == "" {
				issues = append(issues, constants.ErrAIMissingCodeResource)
			} else if utf8.RuneCountInString(resource) > constants.AIMaxResourceLength {
				issues = append(issues, constants.ErrAIResourceTooLong)
			}
		} else {
			resource = ""
		}

		fingerprint := strings.Join(strings.Fields(strings.ToLower(question)), " ")
		if questionMedia == constants.MediaCode {
			fingerprint += "\x00" + strings.Join(strings.Fields(resource), " ")
		}
		if question != "" && seen[fingerprint] {
			issues = append(issues, constants.ErrAIDuplicateQuestion)
		}

		options := make([]string, 0, len(item.Options))
		blankOption := false
		longOption := false
		for _, option := range item.Options {
			trimmed := strings.TrimSpace(option)
			if trimmed == "" {
				blankOption = true
			}
			if utf8.RuneCountInString(trimmed) > constants.AIMaxOptionLength {
				longOption = true
			}
			options = append(options, trimmed)
		}
		if blankOption {
			issues = append(issues, constants.ErrAIBlankOption)
		}
		if longOption {
			issues = append(issues, constants.ErrAIOptionTooLong)
		}

		if len(options) < constants.AIMinOptionsPerQuestion {
			issues = append(issues, constants.ErrInsufficientOptions)
		} else if len(options) > constants.AIMaxOptionsPerQuestion {
			issues = append(issues, fmt.Sprintf(constants.ErrAITooManyOptions, constants.AIMaxOptionsPerQuestion))
		}

		if hasDuplicateOption(options, optionsMedia) {
			issues = append(issues, constants.ErrAIDuplicateOption)
		}

		questionType, questionTypeOK := normalizeAIQuestionType(item.QuestionType)
		if !questionTypeOK {
			issues = append(issues, fmt.Sprintf("%s (got %q)", constants.ErrQuestionType, item.QuestionType))
		}

		correctAnswer := item.CorrectAnswer
		answers := []int{correctAnswer}
		if questionType == constants.Survey {
			correctAnswer = 0
			answers = make([]int, 0, len(options))
			for index := range options {
				answers = append(answers, index+1)
			}
		} else if item.CorrectAnswer < 1 || item.CorrectAnswer > len(options) {
			issues = append(issues, fmt.Sprintf("%s (got %d, %d options)", constants.ErrAICorrectAnswerOutOfRange, item.CorrectAnswer, len(options)))
		}

		if len(issues) > 0 {
			result.Issues = append(result.Issues, fmt.Sprintf("question %d: %s", questionNo, strings.Join(issues, "; ")))
			continue
		}

		optionMap := make(map[string]string, len(options))
		for idx, option := range options {
			optionMap[strconv.Itoa(idx+1)] = option
		}

		seen[fingerprint] = true

		result.Questions = append(result.Questions, models.Question{
			Question:          question,
			Type:              questionType,
			Options:           optionMap,
			Answers:           answers,
			Points:            points,
			DurationInSeconds: opts.Duration,
			QuestionMedia:     questionMedia,
			OptionsMedia:      optionsMedia,
			Resource:          sql.NullString{String: resource, Valid: true},
		})

		result.Accepted = append(result.Accepted, structs.AIQuestion{
			Question:      question,
			QuestionType:  aiQuestionTypeWord(questionType),
			QuestionMedia: questionMedia,
			Resource:      resource,
			Options:       options,
			OptionsMedia:  optionsMedia,
			CorrectAnswer: correctAnswer,
			Explanation:   TruncateRunes(strings.TrimSpace(item.Explanation), constants.AIMaxExplanationLength),
		})
	}

	if truncated {
		result.Issues = append(result.Issues, constants.AIQuestionsTruncated)
	}

	if len(result.Questions) == 0 {
		return result, errors.New(constants.ErrAINoValidQuestions)
	}

	return result, nil
}

func hasDuplicateOption(options []string, optionsMedia string) bool {
	seen := make(map[string]bool, len(options))
	for _, option := range options {
		key := option
		if optionsMedia != constants.MediaCode {
			key = strings.ToLower(option)
		}
		if seen[key] {
			return true
		}
		seen[key] = true
	}
	return false
}

func stripCodeFence(content string) string {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "```") {
		return trimmed
	}

	trimmed = strings.TrimPrefix(trimmed, "```")
	if newline := strings.IndexByte(trimmed, '\n'); newline != -1 {
		if firstLine := strings.TrimSpace(trimmed[:newline]); !strings.Contains(firstLine, "{") && !strings.Contains(firstLine, "[") {
			trimmed = trimmed[newline+1:]
		}
	}

	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(trimmed), "```"))
}

func sliceBetween(content string, opening, closing byte) string {
	start := strings.IndexByte(content, opening)
	end := strings.LastIndexByte(content, closing)
	if start == -1 || end == -1 || end <= start {
		return ""
	}
	return content[start : end+1]
}
