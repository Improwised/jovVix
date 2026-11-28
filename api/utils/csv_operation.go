package utils

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/google/uuid"
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
	// Row is the record's row in the uploaded file (header is row 1), so blank lines ParseCSV
	// skips don't shift later rows. 0 when the row did not come from ParseCSV.
	Row int `csv:"-"`
}

// rowNumber is the spreadsheet row of the i-th question.
func (u Question) rowNumber(i int) int {
	if u.Row > 0 {
		return u.Row
	}
	return i + 2
}

// CSVField is a column the importer understands. Name is also the header in the sample CSV.
type CSVField struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

// CSVFields lists the importable fields in display order. Required matches what row
// validation needs: text, type, two options and the correct answer.
var CSVFields = []CSVField{
	{"Question Text", true}, {"Question Type", true}, {"Points", false},
	{"Option 1", true}, {"Option 2", true}, {"Option 3", false}, {"Option 4", false}, {"Option 5", false},
	{"Correct Answer", true}, {"Question Media", false}, {"Options Media", false}, {"Resource", false},
}

// CSVImport is a parsed upload: its headers, the field -> header mapping that was applied,
// the required fields left unmapped, and the rows in importer fields.
type CSVImport struct {
	Headers          []string          `json:"headers"`
	Mapping          map[string]string `json:"mapping"`
	UnmappedRequired []string          `json:"unmapped_required"`
	Questions        []Question        `json:"-"`
}

// normalizeHeader makes "question_text", "Question Text " and "QUESTION-TEXT" equal.
func normalizeHeader(h string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(h) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ParseCSV reads a CSV with any headers. mapping maps a field name to one of the file's headers;
// when it is empty, each field maps to the header with the same normalized name. A mapped header
// that the file lacks, or an unknown field, is an error. Unmapped required fields are reported,
// not an error, so the caller can ask the user to map them.
func ParseCSV(data []byte, mapping map[string]string) (CSVImport, error) {
	result := CSVImport{Mapping: map[string]string{}, UnmappedRequired: []string{}}

	reader := csv.NewReader(bytes.NewReader(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))))
	reader.FieldsPerRecord = -1 // a short row reads as blank cells
	records, err := reader.ReadAll()
	if err != nil {
		return result, err
	}
	if len(records) < 2 {
		return result, fmt.Errorf(constants.ErrEmptyFile)
	}
	result.Headers = records[0]

	headerIndex := map[string]int{}
	normalizedIndex := map[string]int{}
	for i, h := range result.Headers {
		if _, seen := headerIndex[h]; !seen {
			headerIndex[h] = i
		}
		if _, seen := normalizedIndex[normalizeHeader(h)]; !seen {
			normalizedIndex[normalizeHeader(h)] = i
		}
	}

	fieldIndex := map[string]int{}
	known := map[string]bool{}
	for _, f := range CSVFields {
		known[f.Name] = true
	}
	if len(mapping) == 0 {
		for _, f := range CSVFields {
			if i, ok := normalizedIndex[normalizeHeader(f.Name)]; ok {
				fieldIndex[f.Name] = i
				result.Mapping[f.Name] = result.Headers[i]
			}
		}
	} else {
		for field, header := range mapping {
			if header == "" {
				continue
			}
			if !known[field] {
				return result, fmt.Errorf("%s: unknown field %q", constants.ErrInvalidCSVMapping, field)
			}
			i, ok := headerIndex[header]
			if !ok {
				return result, fmt.Errorf("%s: the file has no column %q", constants.ErrInvalidCSVMapping, header)
			}
			fieldIndex[field] = i
			result.Mapping[field] = header
		}
	}
	for _, f := range CSVFields {
		if _, ok := fieldIndex[f.Name]; f.Required && !ok {
			result.UnmappedRequired = append(result.UnmappedRequired, f.Name)
		}
	}

	cell := func(record []string, field string) string {
		if i, ok := fieldIndex[field]; ok && i < len(record) {
			return record[i]
		}
		return ""
	}
	for j, record := range records[1:] {
		if strings.TrimSpace(strings.Join(record, "")) == "" {
			continue // blank line
		}
		result.Questions = append(result.Questions, Question{
			Question: cell(record, "Question Text"), Type: cell(record, "Question Type"), Points: cell(record, "Points"),
			Option1: cell(record, "Option 1"), Option2: cell(record, "Option 2"), Option3: cell(record, "Option 3"),
			Option4: cell(record, "Option 4"), Option5: cell(record, "Option 5"), CorrectAnswer: cell(record, "Correct Answer"),
			QuestionMedia: cell(record, "Question Media"), OptionsMedia: cell(record, "Options Media"), Resource: cell(record, "Resource"),
			Row: j + 2,
		})
	}
	if len(result.Questions) == 0 {
		return result, fmt.Errorf(constants.ErrEmptyFile)
	}
	return result, nil
}

// ParseCSVFile reads an uploaded file with ParseCSV and requires every required field mapped.
func ParseCSVFile(fileName string, mapping map[string]string) ([]Question, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	parsed, err := ParseCSV(data, mapping)
	if err != nil {
		return nil, err
	}
	if len(parsed.UnmappedRequired) > 0 {
		return nil, fmt.Errorf("%s: %s", constants.ErrUnmappedCSVColumns, strings.Join(parsed.UnmappedRequired, ", "))
	}
	return parsed.Questions, nil
}

// ParseCSVMapping decodes the optional "mapping" form value: a JSON object of field -> header.
func ParseCSVMapping(raw string) (map[string]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	mapping := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &mapping); err != nil {
		return nil, fmt.Errorf("%s: %v", constants.ErrInvalidCSVMapping, err)
	}
	return mapping, nil
}

// ValidateCSVFileFormat reads a CSV whose headers match the sample file.
func ValidateCSVFileFormat(fileName string) ([]Question, error) {
	return ParseCSVFile(fileName, nil)
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
	if duration > constants.MaximumDurationInSeconds {
		duration = constants.MaximumDurationInSeconds
	}

	var validQuestions []models.Question
	var rowErrors []string

	for i, u := range questions {
		rowNo := u.rowNumber(i)

		question, rowIssues := validateCSVRow(u, i+1, duration)
		if len(rowIssues) > 0 {
			rowErrors = append(rowErrors, fmt.Sprintf("row %d: %s", rowNo, strings.Join(rowIssues, "; ")))
			continue
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return validQuestions, err
		}
		question.ID = id
		validQuestions = append(validQuestions, question)
	}

	if len(rowErrors) > 0 {
		return nil, fmt.Errorf("%s %s", constants.ErrInvalidCSVRows, strings.Join(rowErrors, " | "))
	}

	return validQuestions, nil
}

// CSVRowResult is one data row of a preview: its spreadsheet row number, its cells by field
// name, and what is wrong with it.
type CSVRowResult struct {
	Row    int               `json:"row"`
	Values map[string]string `json:"values"`
	Errors []string          `json:"errors"`
}

// PreviewCSVRows validates each row the way ExtractQuestionsFromCSV does, without stopping.
func PreviewCSVRows(questions []Question) []CSVRowResult {
	results := make([]CSVRowResult, 0, len(questions))
	for i, u := range questions {
		_, issues := validateCSVRow(u, i+1, constants.MaximumDurationInSeconds)
		if issues == nil {
			issues = []string{}
		}
		results = append(results, CSVRowResult{
			Row: u.rowNumber(i),
			Values: map[string]string{
				"Question Text": u.Question, "Question Type": u.Type, "Points": u.Points,
				"Option 1": u.Option1, "Option 2": u.Option2, "Option 3": u.Option3, "Option 4": u.Option4,
				"Option 5": u.Option5, "Correct Answer": u.CorrectAnswer, "Question Media": u.QuestionMedia,
				"Options Media": u.OptionsMedia, "Resource": u.Resource,
			},
			Errors: issues,
		})
	}
	return results
}

// validateCSVRow checks one row and builds its question (without an ID).
func validateCSVRow(u Question, order int, duration int) (models.Question, []string) {
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

	// Points: optional, but if present must fall inside the allowed range.
	points := constants.DefaultCsvPoints
	if strings.TrimSpace(u.Points) != "" {
		parsedPoints, convErr := strconv.Atoi(strings.TrimSpace(u.Points))
		if convErr != nil || parsedPoints < constants.MinimumPoints || parsedPoints > constants.MaximumPoints {
			rowIssues = append(rowIssues, fmt.Sprintf("%s (got %q, allowed: %d-%d)", constants.ErrInvalidPoints, u.Points, constants.MinimumPoints, constants.MaximumPoints))
		} else {
			points = parsedPoints
		}
	}

	return models.Question{
		Question:          u.Question,
		Type:              questionType,
		Options:           options,
		Answers:           answers,
		Points:            int16(points),
		DurationInSeconds: duration,
		OrderNumber:       order,
		QuestionMedia:     questionMedia,
		OptionsMedia:      optionsMedia,
		Resource:          sql.NullString{String: u.Resource, Valid: true},
	}, rowIssues
}
