package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
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
	ID          uuid.UUID         `json:"id" db:"id"`
	Question    string            `json:"question" db:"question"`
	Options     map[string]string `json:"options" db:"options"`
	Answers     []int             `json:"answers" db:"answers"`
	Score       int16             `json:"score,omitempty" db:"score"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
	OrderNumber int               `json:"order"`
}

// QuizModel implements quiz related database operations
type QuestionModel struct {
	db *goqu.Database
}

// InitQuizModel initializes the QuizModel
func InitQuestionModel(goquDB *goqu.Database) *QuestionModel {
	return &QuestionModel{db: goquDB}
}

func (model *QuestionModel) CalculateScore(userAnswer structs.ReqAnswerSubmit) (sql.NullInt16, error) {
	var answers []int = []int{}
	var answerScore int16
	var answerBytes []byte = []byte{}
	var score sql.NullInt16 = sql.NullInt16{}

	rows, err := model.db.Select(goqu.I("answers"), goqu.I("score")).From(QuestionTable).Where(goqu.I("id").Eq(userAnswer.QuestionId.String())).Executor().Query()

	if err != nil {
		return score, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&answerBytes, &answerScore)
		if err != nil {
			return score, err
		}
	}

	err = json.Unmarshal(answerBytes, &answers)
	if err != nil {
		return score, err
	}

	// check type of the question
	actualAnswerLen := len(answers)
	userAnswerLen := len(userAnswer.AnswerKeys)

	// if not attempted
	if userAnswerLen == 0 {
		return score, nil
	}

	score.Valid = true
	// for mcq type question
	if actualAnswerLen == 1 {
		if answers[0] == userAnswer.AnswerKeys[0] {
			score.Int16 = answerScore
			return score, nil
		}
		return score, nil
	}

	var noOfMatches int = 0
	for _, actualAnswer := range answers {
		for _, userAnswer := range userAnswer.AnswerKeys {
			if actualAnswer == userAnswer {
				noOfMatches += 1
				if noOfMatches == userAnswerLen {
					break
				}
			}
		}
	}
	score.Int16 = int16(noOfMatches) * answerScore
	return score, nil
}
