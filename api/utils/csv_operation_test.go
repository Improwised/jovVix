package utils

import (
	"os"
	"testing"

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
		assert.NoError(t, err)
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
		assert.NoError(t, err)
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

	t.Run("Empty Options", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Choose the best option",
				Type:          "single answer",
				Points:        "10",
				CorrectAnswer: "1",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "")
		assert.NoError(t, err)
		assert.Len(t, validQuestions, 1)

		assert.Equal(t, map[string]string{}, validQuestions[0].Options)
	})

	t.Run("Different Question Types", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "Select the correct statement",
				Type:          "single answer",
				CorrectAnswer: "1",
			},
			{
				Question:      "Provide your feedback",
				Type:          "survey",
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
				CorrectAnswer: "7",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "120")
		assert.NoError(t, err)
		assert.Len(t, validQuestions, 1)
		assert.Equal(t, 120, validQuestions[0].DurationInSeconds)
	})

	t.Run("Empty Fields", func(t *testing.T) {
		questions := []Question{
			{
				Question:      "",
				Type:          "single answer",
				CorrectAnswer: "",
			},
		}

		validQuestions, err := ExtractQuestionsFromCSV(questions, "20")
		assert.NoError(t, err)
		assert.Len(t, validQuestions, 1)

		assert.Equal(t, "", validQuestions[0].Question)
		assert.Equal(t, []int{}, validQuestions[0].Answers)
	})
}
