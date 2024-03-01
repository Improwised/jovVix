package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// Question model
type QuestionAnswer struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	Question     string         `json:"question" db:"question"`
	Options      string         `json:"options" db:"options"`
	Answers      string         `json:"answers" db:"answers"`
	Score        int            `json:"score,omitempty" db:"score"`
	NextQuestion sql.NullString `json:"next_question" db:"next_question"`
	CreatedAt    time.Time      `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt    time.Time      `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// QuestionModel implements question related database operations
type QuestionAnswerModel struct {
	db *goqu.Database
}

// InitQuestionModel initializes the QuestionModel
func InitQuestionAnswerModel(goqu *goqu.Database) *QuestionAnswerModel {
	return &QuestionAnswerModel{db: goqu}
}
