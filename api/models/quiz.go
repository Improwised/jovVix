package models

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// Quiz model
type Quiz struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" validate:"required"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatorID   string    `json:"creator_id,omitempty" db:"creator_id"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// QuizModel implements quiz related database operations
type QuizModel struct {
	db *goqu.Database
}

// InitQuizModel initializes the QuizModel
func InitQuizModel(goquDB *goqu.Database) *QuizModel {
	return &QuizModel{db: goquDB}
}

func (q *Quiz) GetQuestions() {

}

func (q *Quiz) GetQuizMeta() {

}

func (q *Quiz) IsAdmin() {

}

func (q *Quiz) SubmitAnswer() {

}

func (q *Quiz) RegisterPlayer() {

}

func (q *Quiz) StartQuiz() {

}

func (q *Quiz) SkipQuiz() {

}
