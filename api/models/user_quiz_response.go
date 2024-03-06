package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const UserQuizResponseTable = "user_quiz_response"

// Question model
type UserQuizResponse struct {
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
type UserQuizResponseModel struct {
	db *goqu.Database
}

// InitQuestionModel initializes the QuestionModel
func InitUserQuizResponseModel(goqu *goqu.Database) *UserQuizResponseModel {
	return &UserQuizResponseModel{db: goqu}
}
