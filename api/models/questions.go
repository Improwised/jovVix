package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	QuizId            uuid.UUID         `json:"quiz_id" db:"quiz_id"`
	Question          string            `json:"question" db:"question"`
	Type              int               `json:"type" db:"type"`
	Options           map[string]string `json:"options" db:"options"`
	Answers           []int             `json:"answers" db:"answers,omitempty"`
	Points            int16             `json:"points,omitempty" db:"points,omitempty"`
	DurationInSeconds int               `json:"duration" db:"duration_in_seconds"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at,omitempty"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at,omitempty"`
	OrderNumber       int               `json:"order" db:"order_no"`
	QuestionMedia     string            `json:"question_media" db:"question_media"`
	OptionsMedia      string            `json:"options_media" db:"options_media"`
	Resource          sql.NullString    `json:"resource" db:"resource"`
}

type QuestionForUser struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	Question          string            `json:"question" db:"question"`
	RawOptions        []byte            `json:"omitempty" db:"options"`
	Options           map[string]string `json:"options" db:"omitempty"`
	DurationInSeconds int               `json:"duration" db:"duration_in_seconds"`
	OrderNumber       int               `json:"order" db:"order_no"`
	Points            int               `json:"points" db:"points"`
	QuestionMedia     string            `json:"question_media" db:"question_media"`
	OptionsMedia      string            `json:"options_media" db:"options_media"`
	Resource          sql.NullString    `json:"resource" db:"resource"`
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
			"type":                question.Type,
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
		model.logger.Debug("error in registerQuiz", zap.Error(err))
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

	if err != nil {
		return quizId, err
	}

	if !ok {
		return quizId, sql.ErrNoRows
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
			"type":                question.Type,
			"options":             string(options),
			"answers":             string(answers),
			"points":              question.Points,
			"duration_in_seconds": question.DurationInSeconds,
			"question_media":      question.QuestionMedia,
			"options_media":       question.OptionsMedia,
			"resource":            question.Resource.String,
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

func (model *QuestionModel) GetAnswersPointsDurationType(QuestionID string) ([]int, int16, int, int, error) {

	var answers []int = []int{}
	var answerDurationInSeconds int
	var answerBytes []byte = []byte{}
	var answerPoints int16
	var questionType int

	rows, err := model.db.Select(goqu.I("answers"), goqu.I("points"), goqu.I("duration_in_seconds"), goqu.I("type")).From(QuestionTable).Where(goqu.I("id").Eq(QuestionID)).Executor().Query()

	if err != nil {
		return answers, answerPoints, answerDurationInSeconds, 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&answerBytes, &answerPoints, &answerDurationInSeconds, &questionType)
		if err != nil {
			return answers, answerPoints, answerDurationInSeconds, 0, err
		}
	}

	err = json.Unmarshal(answerBytes, &answers)
	if err != nil {
		return answers, answerPoints, answerDurationInSeconds, 0, err
	}

	return answers, answerPoints, answerDurationInSeconds, questionType, nil
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
			"question_media",
			"options_media",
			"resource",
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

func (model *QuestionModel) GetTotalQuestionCount(activeQuizId string) (int64, error) {
	return model.db.From(ActiveQuizQuestionsTable).Where(goqu.Ex{
		"active_quiz_id": activeQuizId,
	}).Count()
}

func (model *QuestionModel) ListQuestionByQuizId(QuizId string, media string) ([]Question, error) {
	var questions []Question

	query := model.db.From(QuestionTable).
		Join(
			goqu.T(constants.QuizQuestionsTable),
			goqu.On(goqu.I("quiz_questions.question_id").Eq(goqu.I("questions.id"))),
		).
		Where(
			goqu.I("quiz_questions.quiz_id").Eq(QuizId),
		).
		Select(
			"questions.id",
			"questions.question",
			"questions.question_media",
			"questions.options_media",
		)

	if media != "" {
		query = query.Where(goqu.Or(
			goqu.I("questions.question_media").Eq(media),
			goqu.I("questions.options_media").Eq(media),
		))
	}

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	err = model.db.ScanStructs(&questions, sql, args...)

	return questions, err
}

func (model *QuestionModel) UpdateQuestionsResourceById(id, resource string) error {

	result, err := model.db.Update(QuestionTable).Set(goqu.Record{
		"resource":   resource,
		"updated_at": goqu.L("now()"),
	}).Where(goqu.I("id").Eq(id)).Executor().Exec()
	if err != nil {
		return err
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (model *QuestionModel) UpdateQuestionsOptionById(id, keyPath, data string) error {

	jsonValue := fmt.Sprintf("\"%s\"", data)

	result, err := model.db.Update(QuestionTable).Set(goqu.Record{
		"options": goqu.L("jsonb_set(options::jsonb, '{" + keyPath + "}', '" + jsonValue + "')"),
	}).Where(goqu.I("id").Eq(id)).Executor().Exec()
	if err != nil {
		return err
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return sql.ErrNoRows
	}
	return nil
}
