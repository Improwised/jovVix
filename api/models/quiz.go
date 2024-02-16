package models

import "github.com/doug-martin/goqu/v9"

type Quiz struct {
	db *goqu.Database
}

func NewQuiz(db *goqu.Database) *Quiz{
	return &Quiz{db}
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
