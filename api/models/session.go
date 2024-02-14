package models

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// QuizSession model
type QuizSession struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Code          int       `json:"code" db:"code"`
	Title         string    `json:"title,omitempty" db:"title"`
	QuizID        uuid.UUID `json:"quiz_id" db:"quiz_id"`
	AdminID       string    `json:"admin_id,omitempty" db:"admin_id"`
	MaxAttempt    int       `json:"max_attempt" db:"max_attempt"`
	ActivatedTo   time.Time `json:"activated_to,omitempty" db:"activated_to"`
	ActivatedFrom time.Time `json:"activated_from,omitempty" db:"activated_from"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	QuizAnalysis  string    `json:"quiz_analysis,omitempty" db:"quiz_analysis"`
	CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// QuizSessionModel implements quiz session related database operations
type QuizSessionModel struct {
	db *goqu.Database
}

// InitQuizSessionModel initializes the QuizSessionModel
func InitQuizSessionModel(goqu *goqu.Database) *QuizSessionModel {
	return &QuizSessionModel{goqu}
}
