package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const ActiveQuizQuestionsTable = "active_quiz_questions"

type ActiveQuizQuestions struct {
	ID            uuid.UUID `json:"id" db:"id"`
	QuestionID    uuid.UUID `json:"question_id" db:"question_id"`
	NextQuestion  uuid.UUID `json:"next_question" db:"next_question"`
	QuizSessionID uuid.UUID `json:"active_quiz_id" db:"active_quiz_id"`
	OrderNo       int       `json:"order_no" db:"order_no"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

const QuestionTable = "questions"

type Question struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	Question          string            `json:"question" db:"question"`
	Type 			  int 				`json:"type" db:"type"`
	Options           map[string]string `json:"options" db:"options"`
	Answers           []int             `json:"answers" db:"answers,omitempty"`
	Points            int16             `json:"points,omitempty" db:"points,omitempty"`
	DurationInSeconds int               `json:"duration" db:"duration_in_seconds"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at,omitempty"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at,omitempty"`
	OrderNumber       int               `json:"order" db:"order_no"`
}

type QuestionForUser struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	Question          string            `json:"question" db:"question"`
	RawOptions        []byte            `json:"omitempty" db:"options"`
	Options           map[string]string `json:"options" db:"omitempty"`
	DurationInSeconds int               `json:"duration" db:"duration_in_seconds"`
	OrderNumber       int               `json:"order" db:"order_no"`
	Points            int               `json:"points" db:"points"`
}

// QuizModel implements quiz related database operations
type QuestionModel struct {
	db     *goqu.Database
	logger *zap.Logger
}

// InitQuizModel initializes the QuizModel
func InitQuestionModel(goquDB *goqu.Database, logger *zap.Logger) *QuestionModel {
	return &QuestionModel{db: goquDB, logger: logger}
}

func (model *QuestionModel) CreateQuestions(quizId uuid.UUID, questions []Question) ([]uuid.UUID, error) {
	ids := []uuid.UUID{}
	records := []goqu.Record{}

	for _, question := range questions {
		options, err := json.Marshal(question.Options)
		if err != nil {
			return ids, err
		}

		answers, err := json.Marshal(question.Answers)
		if err != nil {
			return ids, err
		}

		records = append(records, goqu.Record{
			"id":                  question.ID,
			"question":            question.Question,
			"type": 			   question.Type,
			"options":             string(options),
			"answers":             string(answers),
			"points":              question.Points,
			"duration_in_seconds": question.DurationInSeconds,
		})
	}

	err := model.db.Insert(QuestionTable).Rows(
		records,
	).Returning("id").Executor().ScanVals(&ids)

	if err != nil {
		return ids, err
	}

	return ids, nil
}

func (model *QuestionModel) RegisterQuestions(userId string, title string, description string, questions []Question) (uuid.UUID, error) {

	isOk := false
	transaction, err := model.db.Begin()

	if err != nil {
		return uuid.UUID{}, err
	}

	defer func() {
		if isOk {
			err := transaction.Commit()
			if err != nil {
				model.logger.Error("error during commit in register question", zap.Error(err))
			}
		} else {
			err := transaction.Rollback()
			if err != nil {
				model.logger.Error("error during rollback in register question", zap.Error(err))
			}
		}
	}()

	quizId, err := registerQuiz(transaction, title, description, userId)

	if err != nil {
		return quizId, err
	}

	ids, err := registerQuestions(transaction, questions)

	if err != nil {
		return quizId, err
	}

	err = registerQuestionToQuizzes(transaction, quizId, ids)

	if err != nil {
		return quizId, err
	}

	isOk = true
	return quizId, nil
}

func registerQuiz(transaction *goqu.TxDatabase, title, description, userId string) (uuid.UUID, error) {
	quizId, err := uuid.NewUUID()

	if err != nil {
		return quizId, err
	}

	ok, err := transaction.Insert(QuizzesTable).Rows(
		goqu.Record{
			"id":          quizId,
			"title":       title,
			"description": sql.NullString{Valid: description != "", String: description},
			"creator_id":  userId,
		},
	).Returning("id").Executor().ScanVal(&quizId)

	if !ok {
		return quizId, sql.ErrNoRows
	}

	if err != nil {
		return quizId, err
	}

	return quizId, nil
}

func registerQuestions(transaction *goqu.TxDatabase, questions []Question) ([]uuid.UUID, error) {
	ids := []uuid.UUID{}
	records := []goqu.Record{}

	for _, question := range questions {
		options, err := json.Marshal(question.Options)
		if err != nil {
			return ids, err
		}

		answers, err := json.Marshal(question.Answers)
		if err != nil {
			return ids, err
		}

		records = append(records, goqu.Record{
			"id":                  question.ID,
			"question":            question.Question,
			"type": 			   question.Type,
			"options":             string(options),
			"answers":             string(answers),
			"points":              question.Points,
			"duration_in_seconds": question.DurationInSeconds,
		})
	}

	err := transaction.Insert(QuestionTable).Rows(
		records,
	).Returning("id").Executor().ScanVals(&ids)

	if err != nil {
		return ids, err
	}

	return ids, err
}

func registerQuestionToQuizzes(transaction *goqu.TxDatabase, quizId uuid.UUID, questionIds []uuid.UUID) error {
	records := []goqu.Record{}
	for questionIdIndex, questionId := range questionIds {

		id, err := uuid.NewUUID()

		if err != nil {
			return err
		}

		nextQuestion := uuid.NullUUID{}
		if questionIdIndex+1 != len(questionIds) {
			nextQuestion.Valid = true
			nextQuestion.UUID = questionIds[questionIdIndex+1]
		}

		records = append(records,
			goqu.Record{
				"id":            id,
				"question_id":   questionId,
				"quiz_id":       quizId,
				"next_question": nextQuestion,
			},
		)
	}

	rows, err := transaction.Insert("quiz_questions").Rows(
		records,
	).Executor().Exec()

	if err != nil {
		return err
	}

	insertedRowCount, err := rows.RowsAffected()

	if err != nil {
		return err
	}

	if insertedRowCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (model *QuestionModel) GetAnswersPointsDurationOptions(QuestionID string) ([]int, int16, int, map[string]string, error) {

	var answers []int = []int{}
	var options map[string]string
	var answerDurationInSeconds int
	var answerBytes []byte = []byte{}
	var answerPoints int16
	var optionsBytes []byte

	rows, err := model.db.Select(goqu.I("answers"), goqu.I("points"), goqu.I("duration_in_seconds"), goqu.I("options")).From(QuestionTable).Where(goqu.I("id").Eq(QuestionID)).Executor().Query()

	if err != nil {
		return answers, answerPoints, answerDurationInSeconds, options, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&answerBytes, &answerPoints, &answerDurationInSeconds, &optionsBytes)
		if err != nil {
			return answers, answerPoints, answerDurationInSeconds, options, err
		}
	}

	err = json.Unmarshal(answerBytes, &answers)
	if err != nil {
		return answers, answerPoints, answerDurationInSeconds, options, err
	}

	err = json.Unmarshal(optionsBytes, &options)
	if err != nil {
		return answers, answerPoints, answerDurationInSeconds, options, err
	}

	return answers, answerPoints, answerDurationInSeconds, options, nil
}

func (model *QuestionModel) GetCurrentQuestion(id uuid.UUID) (QuestionForUser, error) {
	var question QuestionForUser

	ok, err := model.db.From(constants.QuestionsTable).
		Select(
			goqu.I(constants.QuestionsTable+".id"),
			"order_no",
			"duration_in_seconds",
			"question",
			"options",
			"points",
		).InnerJoin(
		goqu.T(constants.ActiveQuizQuestionsTable), goqu.On(goqu.I(constants.QuestionsTable+".id").Eq(goqu.I(constants.ActiveQuizQuestionsTable+".question_id")))).
		Where(goqu.Ex{
			constants.QuestionsTable + ".id": id.String(),
		}).Limit(1).ScanStruct(&question)

	if !ok && err == nil {
		return question, sql.ErrNoRows
	} else if !ok || err != nil {
		return question, err
	} else {
		err = json.Unmarshal(question.RawOptions, &question.Options)
		if err != nil {
			return QuestionForUser{}, err
		}
		return question, nil
	}
}
