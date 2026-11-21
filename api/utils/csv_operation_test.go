package utils

import (
	"os"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/stretchr/testify/assert"
)

func createTempCSV(content string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "test-*.csv")
	if err != nil {
		return nil, err
	}
	_, err = tempFile.Write([]byte(content))
	if err != nil {
		return nil, err
	}
	err = tempFile.Close()
	return tempFile, err
}

func TestValidateCSVFileFormat(t *testing.T) {

	t.Run("Valid CSV File", func(t *testing.T) {

		csvContent := `Question Text,Question Type,Points,Option 1,Option 2,Option 3,Option 4,Option 5,Correct Answer,Question Media,Options Media,Resource
"Sample Question 1","MCQ","10","Option A","Option B","Option C","Option D","","Option A","","","Resource 1"
"Sample Question 2","MCQ","5","Option A","Option B","Option C","Option D","","Option B","","","Resource 2"`

		tempFile, err := createTempCSV(csvContent)
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		questions, err := ValidateCSVFileFormat(tempFile.Name())
		assert.NoError(t, err)
		assert.Len(t, questions, 2)
		assert.Equal(t, "Sample Question 1", questions[0].Question)
		assert.Equal(t, "MCQ", questions[0].Type)
		assert.Equal(t, "10", questions[0].Points)
		assert.Equal(t, "Option A", questions[0].CorrectAnswer)
	})

	t.Run("Invalid CSV File Format (missing headers)", func(t *testing.T) {
		csvContent := `Sample Question 1,MCQ,10,Option A,Option B,Option C,Option D,,Option A,,,Resource 1`

		tempFile, err := createTempCSV(csvContent)
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		questions, err := ValidateCSVFileFormat(tempFile.Name())
		assert.Error(t, err)
		assert.Empty(t, questions)

	})

	t.Run("Non-Existent File", func(t *testing.T) {
		_, err := ValidateCSVFileFormat("non-existent-file.csv")
		assert.Error(t, err)
	})

	t.Run("Empty CSV File", func(t *testing.T) {
		tempFile, err := createTempCSV("")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())

		questions, err := ValidateCSVFileFormat(tempFile.Name())
		assert.Error(t, err)
		assert.Empty(t, questions)
	})
}

func TestExtractQuestionsFromCSV(t *testing.T) {
	t.Run("Valid Data", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "What is the capital of France?",
				Type:          "single answer",
				Points:        "5",
				Option1:       "Paris",
				Option2:       "London",
				Option3:       "Berlin",
				Option4:       "Madrid",
				CorrectAnswer: "1",
				QuestionMedia: "text",
				OptionsMedia:  "text",
				Resource:      "",
			},
			{
				Question:      "Rate the service quality",
				Type:          "survey",
				Points:        "3",
				Option1:       "Poor",
				Option2:       "Average",
				Option3:       "Good",
				Option4:       "Excellent",
				CorrectAnswer: "3|4",
				QuestionMedia: "",
				OptionsMedia:  "text",
				Resource:      "SurveyMonkey",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "45")
		assert.NoError(t, err)
		assert.Len(t, validQuestions, 2)

		// Check first question
		assert.Equal(t, "What is the capital of France?", validQuestions[0].Question)
		assert.Equal(t, 1, validQuestions[0].Type)
		assert.Equal(t, 5, int(validQuestions[0].Points))
		assert.Equal(t, 45, validQuestions[0].DurationInSeconds)
		assert.Equal(t, map[string]string{"1": "Paris", "2": "London", "3": "Berlin", "4": "Madrid"}, validQuestions[0].Options)
		assert.Equal(t, []int{1}, validQuestions[0].Answers)
		assert.Equal(t, "text", validQuestions[0].QuestionMedia)

		// Check second question
		assert.Equal(t, "Rate the service quality", validQuestions[1].Question)
		assert.Equal(t, 2, validQuestions[1].Type)
		assert.Equal(t, 3, int(validQuestions[1].Points))
		assert.Equal(t, 45, validQuestions[1].DurationInSeconds)
		assert.Equal(t, []int{3, 4}, validQuestions[1].Answers)
		assert.Equal(t, "text", validQuestions[1].OptionsMedia)
	})

	t.Run("Different Question Types", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Select the correct statement",
				Type:          "single answer",
				Option1:       "True",
				Option2:       "False",
				CorrectAnswer: "1",
			},
			{
				Question:      "Provide your feedback",
				Type:          "survey",
				Option1:       "Yes",
				Option2:       "No",
				CorrectAnswer: "2",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "60")
		assert.NoError(t, err)
		assert.Len(t, validQuestions, 2)

		// Check question types
		assert.Equal(t, 1, validQuestions[0].Type)
		assert.Equal(t, 2, validQuestions[1].Type)
	})

	t.Run("Custom Time Limit Parsing", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "How many continents are there?",
				Type:          "single answer",
				Points:        "2",
				Option1:       "6",
				Option2:       "7",
				CorrectAnswer: "2",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "120")
		assert.NoError(t, err)
		assert.Len(t, validQuestions, 1)
		assert.Equal(t, 120, validQuestions[0].DurationInSeconds)
	})

	t.Run("Invalid question type is rejected", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "What does CPU stand for?",
				Type:          "MCQ",
				Option1:       "Central Processing Unit",
				Option2:       "Computer Processing Unit",
				CorrectAnswer: "1",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "30")
		assert.Error(t, err)
		assert.Empty(t, validQuestions)
		assert.Contains(t, err.Error(), constants.ErrQuestionType)
	})

	t.Run("Insufficient options is rejected", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Choose the best option",
				Type:          "single answer",
				Points:        "10",
				Option1:       "Only one",
				CorrectAnswer: "1",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "30")
		assert.Error(t, err)
		assert.Empty(t, validQuestions)
		assert.Contains(t, err.Error(), constants.ErrInsufficientOptions)
	})

	t.Run("Missing or invalid time limit config is rejected", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Pick one",
				Type:          "single answer",
				Option1:       "A",
				Option2:       "B",
				CorrectAnswer: "1",
			},
		}

		for _, badLimit := range []string{"", "abc", "0", "-5"} {
			validQuestions, err := ExtractQuestionsFromCSV(questions, badLimit)
			assert.Error(t, err)
			assert.Empty(t, validQuestions)
			assert.Contains(t, err.Error(), constants.ErrInvalidQuestionTimeLimit)
		}
	})

	t.Run("Correct answer referencing missing option is rejected", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "How many continents are there?",
				Type:          "single answer",
				Option1:       "6",
				Option2:       "7",
				CorrectAnswer: "7",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "120")
		assert.Error(t, err)
		assert.Empty(t, validQuestions)
		assert.Contains(t, err.Error(), constants.ErrInvalidCorrectAnswer)
	})

	t.Run("Single answer with multiple correct answers is rejected", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Pick one",
				Type:          "single answer",
				Option1:       "A",
				Option2:       "B",
				CorrectAnswer: "1|2",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "30")
		assert.Error(t, err)
		assert.Empty(t, validQuestions)
		assert.Contains(t, err.Error(), constants.ErrSingleAnswerLength)
	})

	t.Run("Empty fields are rejected", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "",
				Type:          "single answer",
				CorrectAnswer: "",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "20")
		assert.Error(t, err)
		assert.Empty(t, validQuestions)
		assert.Contains(t, err.Error(), constants.ErrEmptyQuestionText)
	})

	t.Run("Invalid media type is rejected", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Pick one",
				Type:          "single answer",
				Option1:       "A",
				Option2:       "B",
				CorrectAnswer: "1",
				QuestionMedia: "video",
				OptionsMedia:  "text",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "30")
		assert.Error(t, err)
		assert.Empty(t, validQuestions)
		assert.Contains(t, err.Error(), constants.ErrInvalidQuestionMedia)
	})

	t.Run("Media is normalized and empty defaults to text", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Pick one",
				Type:          "single answer",
				Option1:       "A",
				Option2:       "B",
				CorrectAnswer: "1",
				QuestionMedia: "  Image ",
				OptionsMedia:  "",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "30")
		assert.NoError(t, err)
		assert.Len(t, validQuestions, 1)
		assert.Equal(t, constants.MediaImage, validQuestions[0].QuestionMedia)
		assert.Equal(t, constants.MediaText, validQuestions[0].OptionsMedia)
	})

	t.Run("Multiple bad rows are all reported", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Good one",
				Type:          "single answer",
				Option1:       "A",
				Option2:       "B",
				CorrectAnswer: "1",
			},
			{
				Question:      "Bad type",
				Type:          "MCQ",
				Option1:       "A",
				Option2:       "B",
				CorrectAnswer: "1",
			},
			{
				Question:      "Bad answer",
				Type:          "single answer",
				Option1:       "A",
				Option2:       "B",
				CorrectAnswer: "9",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "30")
		assert.Error(t, err)
		assert.Empty(t, validQuestions)
		// Row numbers are 1-based with the header as row 1.
		assert.Contains(t, err.Error(), "row 3")
		assert.Contains(t, err.Error(), "row 4")
	})
}

func TestParseCSVMapping(t *testing.T) {
	t.Run("auto-matches headers loosely and strips a BOM", func(t *testing.T) {
		data := "\xef\xbb\xbfquestion_text,QUESTION TYPE,option-1,Option 2,correct answer\nQ1,single answer,A,B,1\n"
		parsed, err := ParseCSV([]byte(data), nil)
		assert.NoError(t, err)
		assert.Empty(t, parsed.UnmappedRequired)
		assert.Equal(t, "question_text", parsed.Mapping["Question Text"])
		assert.Equal(t, "Q1", parsed.Questions[0].Question)
		assert.Equal(t, "1", parsed.Questions[0].CorrectAnswer)
	})

	t.Run("applies a user mapping and reports unmapped required fields", func(t *testing.T) {
		data := "Prompt,Kind,First,Second,Right,Notes\nQ1,survey,A,B,1|2,n\n"
		parsed, err := ParseCSV([]byte(data), nil)
		assert.NoError(t, err)
		assert.Equal(t, []string{"Question Text", "Question Type", "Option 1", "Option 2", "Correct Answer"}, parsed.UnmappedRequired)

		parsed, err = ParseCSV([]byte(data), map[string]string{
			"Question Text": "Prompt", "Question Type": "Kind", "Option 1": "First", "Option 2": "Second",
			"Correct Answer": "Right", "Resource": "Notes", "Points": "",
		})
		assert.NoError(t, err)
		assert.Empty(t, parsed.UnmappedRequired)
		assert.Equal(t, Question{Question: "Q1", Type: "survey", Option1: "A", Option2: "B", CorrectAnswer: "1|2", Resource: "n", Row: 2}, parsed.Questions[0])
	})

	t.Run("rejects a mapping to a missing column or unknown field", func(t *testing.T) {
		data := "A,B\n1,2\n"
		_, err := ParseCSV([]byte(data), map[string]string{"Question Text": "Nope"})
		assert.ErrorContains(t, err, constants.ErrInvalidCSVMapping)
		_, err = ParseCSV([]byte(data), map[string]string{"Bogus": "A"})
		assert.ErrorContains(t, err, constants.ErrInvalidCSVMapping)
	})

	t.Run("short rows read as blank cells and blank lines are skipped", func(t *testing.T) {
		data := "Question Text,Question Type,Option 1,Option 2,Correct Answer,Question Media\nQ1,single answer,A,B,1\n,,,,,\n"
		parsed, err := ParseCSV([]byte(data), nil)
		assert.NoError(t, err)
		assert.Len(t, parsed.Questions, 1)
		assert.Equal(t, "", parsed.Questions[0].QuestionMedia)
	})

	t.Run("preview validates every row without stopping", func(t *testing.T) {
		data := "Question Text,Question Type,Option 1,Option 2,Correct Answer,Question Media\n" +
			"Q1,single answer,A,B,1,\n" +
			",single answer,A,B,3,video\n"
		parsed, err := ParseCSV([]byte(data), nil)
		assert.NoError(t, err)
		rows := PreviewCSVRows(parsed.Questions)
		assert.Len(t, rows, 2)
		assert.Equal(t, 2, rows[0].Row)
		assert.Empty(t, rows[0].Errors)
		assert.Equal(t, 3, rows[1].Row)
		assert.Len(t, rows[1].Errors, 4) // empty text, option 3 missing, so no single answer, bad media
		assert.Equal(t, "video", rows[1].Values["Question Media"])
	})

	t.Run("ParseCSVFile refuses unmapped required fields", func(t *testing.T) {
		tempFile, err := createTempCSV("Prompt,Kind\nQ1,survey\n")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())
		_, err = ParseCSVFile(tempFile.Name(), nil)
		assert.ErrorContains(t, err, constants.ErrUnmappedCSVColumns)
	})
}

// Blank lines are skipped, but later rows keep their row number in the file (header is row 1).
func TestCSVRowNumbersSurviveBlankLines(t *testing.T) {
	data := "Question Text,Question Type,Option 1,Option 2,Correct Answer\n" +
		"Good,single answer,A,B,1\n" +
		",,,,\n" +
		",,,,\n" +
		"Bad answer,single answer,A,B,9\n"
	parsed, err := ParseCSV([]byte(data), nil)
	assert.NoError(t, err)

	preview := PreviewCSVRows(parsed.Questions)
	assert.Len(t, preview, 2)
	assert.Equal(t, 2, preview[0].Row)
	assert.Equal(t, 5, preview[1].Row)
	assert.NotEmpty(t, preview[1].Errors)

	_, err = ExtractQuestionsFromCSV(parsed.Questions, "30")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "row 5:")
	assert.NotContains(t, err.Error(), "row 3:")
}
