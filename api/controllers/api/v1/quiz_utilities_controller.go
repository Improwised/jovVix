package v1

import (
	"database/sql"
	"fmt"
	"io"
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
	"github.com/jszwec/csvutil"
	"go.uber.org/zap"
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

	return questions, nil
}

func ExtractQuestionsFromCSV(questions []Question) ([]models.Question, error) {
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

		var duration int
		durationFromEnv := os.Getenv("QUESTION_TIME_LIMIT")
		if durationFromEnv == "" {
			duration = 30
		} else {
			duration, err = strconv.Atoi(durationFromEnv)
			if err != nil {
				duration = 30
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
		})
	}
	return validQuestions, nil
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

	questions, err := ValidateCSVFileFormat(filePath)
	if err != nil {
		ctrl.logger.Error("file validation failed", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validQuestions, err := ExtractQuestionsFromCSV(questions)
	if err != nil {
		ctrl.logger.Error("file validation failed", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrParsingFile)
	}

	quizId, err := ctrl.helper.QuestionModel.RegisterQuestions(userID, quizTitle, quizDescription, validQuestions)
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
