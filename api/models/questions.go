package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionQuestion struct {
	ID            uuid.UUID `json:"id" db:"id"`
	QuestionID    uuid.UUID `json:"question_id" db:"question_id"`
	NextQuestion  uuid.UUID `json:"next_question" db:"next_question"`
	QuizSessionID uuid.UUID `json:"active_quiz_id" db:"active_quiz_id"`
	OrderNo       int       `json:"order_no" db:"order_no"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Question struct {
	ID          uuid.UUID         `json:"id" db:"id"`
	Question    string            `json:"question" db:"question"`
	Options     map[string]string `json:"options" db:"options"`
	Answers     []string          `json:"answers" db:"answers"`
	Score       int               `json:"score,omitempty" db:"score"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
	OrderNumber int               `json:"order"`
}
