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
		"Points",
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

func ExtractQuestionsFromCSV(fileName string, logger *zap.Logger) ([]models.Question, error) {
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
			return questions, err
		}

		/*
			Index: column
			0:  "Question Text",
			1:  "Question Type",
			2:  "Points",
			3:  "Option 1",
			4:  "Option 2",
			5:  "Option 3",
			6:  "Option 4",
			7:  "Option 5",
			8:  "Correct Answer",

		*/

		// extract type
		var questionType int

		typeInt, err := quizUtilsHelper.CheckQuestionType(strings.ToLower(row[1]))
		if err != nil {
			return questions, err
		} else {
			questionType = typeInt
		}

		var options map[string]string = map[string]string{}
		var optionKey = 1

		// extract options
		for _, option := range []string{row[3], row[4], row[5], row[6], row[7]} {
			if option != "" {
				options[quizUtilsHelper.GetString(optionKey)] = option
				optionKey += 1
			}
		}

		var answers []int

		// extract answers
		for _, answer := range strings.Split(row[8], "|") {
			if _, ok := options[answer]; !ok {
				return questions, fmt.Errorf(fmt.Sprintf("answer not in option list options: %s, answers: %s", options, row[8]))
			}

			answerInt, err := strconv.Atoi(answer)

			if err != nil {
				return questions, fmt.Errorf(fmt.Sprintf("answer string to int fail options: %s, answers: %s", options, row[8]))
			}

			answers = append(answers, answerInt)
		}

		// survey question should have all the options as correct answers
		if questionType == constants.Survey && len(options) != len(answers) {
			return questions, fmt.Errorf(constants.ErrSurveyAnswerLength)
		}

		// single answer question should have single option correct only
		if questionType == constants.SingleAnswer && len(answers) > 1 {
			return questions, fmt.Errorf(constants.ErrSingleAnswerLength)
		}

		// extract score
		var points int16
		if row[2] == "" {
			if questionType == constants.Survey {
				points = 0
			} else {
				points = 1
			}
		} else {
			pointsInt, err := strconv.Atoi(row[2])

			if err != nil {
				return questions, fmt.Errorf(fmt.Sprintf("score string to int fail score: %s", row[2]))
			}

			maximumPointsFromENV := os.Getenv("MAXIMUM_POINTS_PER_QUESTION")
			var maximumPoints int
			if maximumPointsFromENV == "" {
				logger.Debug("MAXIMUM_POINTS_PER_QUESTION env variable is not set, hence taking it from the constants")
				maximumPoints = constants.MaximumPoints
			} else {
				maximumPoints, err = strconv.Atoi(maximumPointsFromENV)
			}
			if err != nil {
				return questions, fmt.Errorf(constants.ErrorTypeConversion)
			}

			minimumPointsFromENV := os.Getenv("MINIMUM_POINTS_PER_QUESTION")
			var minimumPoints int
			if minimumPointsFromENV == "" {
				logger.Debug("MINIMUM_POINTS_PER_QUESTION env variable is not set, hence taking it from the constants")
				minimumPoints = constants.MinimumPoints
			} else {
				minimumPoints, err = strconv.Atoi(minimumPointsFromENV)
			}
			if err != nil {
				return questions, fmt.Errorf(constants.ErrorTypeConversion)
			}

			if pointsInt > maximumPoints {
				return questions, fmt.Errorf(fmt.Sprintf("the points per question should be less than or equal to %d", maximumPoints))
			}

			if pointsInt < minimumPoints {
				return questions, fmt.Errorf(fmt.Sprintf("the points per question should be greater than or equal to %d", minimumPoints))
			}

			points = int16(pointsInt)
		}

		// provide duration
		var duration int
		durationFromEnv := os.Getenv("QUESTION_TIME_LIMIT")
		if durationFromEnv == "" {
			duration = 30
		} else {
			duration, err = strconv.Atoi(durationFromEnv)
			if err != nil {
				duration = 30
				logger.Error("Took default time of 30 seconds per question as convertion error from env", zap.Error(err))
			}
		}

		questions = append(questions,
			models.Question{
				ID:                id,
				Question:          row[0],
				Type:              questionType,
				Options:           options,
				Answers:           answers,
				Points:            points,
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

	questions, err := ExtractQuestionsFromCSV(filePath, ctrl.logger)

	if err != nil {
		if err.Error() == constants.ErrRowsReachesToMaxCount {
			ctrl.logger.Error("file validation failed", zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, err.Error())
		}

		if os.Getenv("MAXIMUM_POINTS_PER_QUESTION") == "" && err.Error() == fmt.Sprintf("the points per question should be less than or equal to %d", constants.MaximumPoints) {
			ctrl.logger.Error("file validation failed", zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, err.Error())
		}

		if os.Getenv("MINIMUM_POINTS_PER_QUESTION") == "" && err.Error() == fmt.Sprintf("the points per question should be greater than or equal to %d", constants.MinimumPoints) {
			ctrl.logger.Error("file validation failed", zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, err.Error())
		}

		if err.Error() == fmt.Sprintf("the points per question should be less than or equal to %s", os.Getenv("MAXIMUM_POINTS_PER_QUESTION")) {
			ctrl.logger.Error("file validation failed", zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, err.Error())
		}

		if err.Error() == fmt.Sprintf("the points per question should be greater than or equal to %s", os.Getenv("MINIMUM_POINTS_PER_QUESTION")) {
			ctrl.logger.Error("file validation failed", zap.Error(err))
			return utils.JSONFail(c, http.StatusBadRequest, err.Error())
		}

		if err.Error() == constants.ErrSurveyAnswerLength {
			ctrl.logger.Error("correct answers are lesser than the number of options in survey question")
			return utils.JSONFail(c, http.StatusBadRequest, err.Error())
		}

		if err.Error() == constants.ErrSingleAnswerLength {
			ctrl.logger.Error(("there are multiple correct answers in single answer type question"))
			return utils.JSONFail(c, http.StatusBadRequest, err.Error())
		}

		if err.Error() == constants.ErrQuestionType {
			ctrl.logger.Error("incorrect question type found in the CSV, upload terminated")
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
