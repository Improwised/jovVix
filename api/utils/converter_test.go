package utils

import (
	"database/sql"
	"testing"

	"github.com/Improwised/jovvix/api/models"
)

func TestResponseToPdfData(t *testing.T) {
	quizReport := []models.AnalyticsBoardAdmin{
		{
			OrderNo:        1,
			UserName:       "John Doe",
			Options:        map[string]string{"0": "Paris", "1": "London", "2": "Berlin", "3": "Rome"},
			CorrectAnswer:  "[1]",
			SelectedAnswer: sql.NullString{String: "[1]", Valid: true},
			Question:       "What is the capital of France?",
			QuestionType:   "single_answer",
		},
		{
			OrderNo:        2,
			UserName:       "John Doe",
			Options:        map[string]string{"0": "Mercury", "1": "Venus", "2": "Earth", "3": "Mars"},
			CorrectAnswer:  "[2]",
			SelectedAnswer: sql.NullString{String: "[2]", Valid: true},
			Question:       "Which planet do we live on?",
			QuestionType:   "single_answer",
		},
	}

	result, order := ResponseToPdfData(quizReport)

	if len(result) != 2 {
		t.Errorf("expected 2 orders, got %d", len(result))
	}
	if order[0] != 1 || order[1] != 2 {
		t.Errorf("unexpected order: %v", order)
	}
	if result[1][0].UserName != "John Doe" {
		t.Errorf("expected John Doe, got %s", result[1][0].UserName)
	}
}
