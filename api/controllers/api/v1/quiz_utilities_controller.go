package v1

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func ValidateCSVFileFormat(fileName string) error {
	// Open the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the first row from the CSV file
	row, err := reader.Read()
	if err != nil {
		return err
	}

	// Required headers
	requiredHeaders := []string{
		"Question Text",
		"Question Type",
		"Time in seconds",
		"Score",
		"Option 1",
		"Option 2",
		"Option 3",
		"Option 4",
		"Option 5",
		"Correct Answer",
	}

	// Validate the headers
	if len(row) < len(requiredHeaders) {
		return fmt.Errorf("CSV file has fewer columns than expected")
	}

	for index, col := range requiredHeaders {
		if col != row[index] {
			return fmt.Errorf("column mismatch at index %d: expected '%s', found '%s'", index, col, row[index])
		}
	}

	return nil
}

func ExtractQuestionsFromCSV(fileName string) ([]models.Question, error) {
	questions := []models.Question{}

	// Open the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		return questions, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the rows from the CSV file
	rows, err := reader.ReadAll()
	if err != nil {
		return questions, err
	}

	for rowNumber, row := range rows[1:] {

		if rowNumber == constants.MaxRows {
			return questions, fmt.Errorf(constants.ErrRowsReachesToMaxCount)
		}

		id, err := uuid.NewUUID()

		if err != nil {
			return questions, nil
		}

		/*
			Index: column
			0:  "Question Text",
			1:  "Question Type",
			2:  "Time in seconds",
			3:  "Score",
			4:  "Option 1",
			5:  "Option 2",
			6:  "Option 3",
			7:  "Option 4",
			8:  "Option 5",
			9:  "Correct Answer",

		*/

		var options map[string]string = map[string]string{}
		var optionKey = 1

		// extract options
		for _, option := range []string{row[4], row[5], row[6], row[7], row[8]} {
			if option != "" {
				options[quizUtilsHelper.GetString(optionKey)] = option
				optionKey += 1
			}
		}

		var answers []int

		// extract answers
		for _, answer := range strings.Split(row[9], "|") {
			if _, ok := options[answer]; !ok {
				return questions, fmt.Errorf(fmt.Sprintf("answer not in option list options: %s, answers: %s", options, row[9]))
			}

			answerInt, err := strconv.Atoi(answer)

			if err != nil {
				return questions, fmt.Errorf(fmt.Sprintf("answer string to int fail options: %s, answers: %s", options, row[9]))
			}

			answers = append(answers, answerInt)
		}

		// extract score
		var score int16
		if row[3] == "" {
			score = 1
		} else {
			scoreInt, err := strconv.Atoi(row[3])

			if err != nil {
				return questions, fmt.Errorf(fmt.Sprintf("score string to int fail score: %s", row[3]))
			}

			score = int16(scoreInt)
		}

		// extract duration
		var duration int
		if row[2] == "" {
			duration = 30
		} else {
			duration, err = strconv.Atoi(row[2])

			if err != nil {
				return questions, fmt.Errorf(fmt.Sprintf("duration string to int fail score: %s", row[2]))
			}
		}

		questions = append(questions,
			models.Question{
				ID:                id,
				Question:          row[0],
				Options:           options,
				Answers:           answers,
				Score:             score,
				DurationInSeconds: duration,
				OrderNumber:       rowNumber + 1,
			},
		)
	}
	return questions, nil
}

func (ctrl *QuizController) CreateQuizByCsv(c *fiber.Ctx) error {

	quizTitle := c.Params(constants.QuizTitle)
	quizDescription := c.FormValue("description")

	if quizTitle == "" {
		ctrl.logger.Error("quiz-title not found")
		return utils.JSONSuccess(c, http.StatusBadRequest, constants.QuizTitleRequired)
	}

	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	filePath := quizUtilsHelper.GetString(c.Locals(constants.FileName))

	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			ctrl.logger.Error("error in deleting file", zap.Error(err))
			return
		}
	}()

	err := ValidateCSVFileFormat(filePath)
	if err != nil {
		ctrl.logger.Error("file validation failed", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	questions, err := ExtractQuestionsFromCSV(filePath)

	if err != nil {
		if err.Error() == constants.ErrRowsReachesToMaxCount {
			ctrl.logger.Error("file validation failed", zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, err.Error())
		}

		ctrl.logger.Error("file validation failed", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrParsingFile)
	}

	quizId, err := ctrl.helper.QuestionModel.RegisterQuestions(userID, quizTitle, quizDescription, questions)

	if err != nil {
		ctrl.logger.Error("error in creating quiz", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrRegisterQuiz)
	}

	return utils.JSONSuccess(c, http.StatusAccepted, quizId)
}

func (ctrl *QuizController) GenerateDemoSession(c *fiber.Ctx) error {
	quizId := c.Params(constants.QuizId)
	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	sessionId, err := ctrl.helper.ActiveQuizModel.CreateActiveQuiz("demo session", quizId, userId, sql.NullTime{}, sql.NullTime{})

	if err != nil {
		ctrl.logger.Error("error in creating demo session", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrCreatingDemoQuiz)
	}

	err = ctrl.helper.ActiveQuizModel.GetQuestionsCopy(sessionId, quizId)
	if err != nil {
		ctrl.logger.Error("error in creating demo session questions", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrCreatingDemoQuiz)
	}

	return utils.JSONSuccess(c, http.StatusAccepted, sessionId)
}
