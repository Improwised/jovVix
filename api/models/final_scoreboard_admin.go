package models

import (
	"github.com/doug-martin/goqu/v9"
)

type FinalScoreBoardAdmin struct {
	UserName string `db:"username" json:"userName"`
	Score    int    `db:"score" json:"score"`
}

type FinalScoreBoardAdminModel struct {
	db *goqu.Database
}

func InitFinalScoreBoardAdminModel(goqu *goqu.Database) (FinalScoreBoardAdminModel, error) {
	return FinalScoreBoardAdminModel{
		db: goqu,
	}, nil
}

// GetScore to send final score after quiz over

func (model *FinalScoreBoardAdminModel) GetScoreForAdmin(activeQuizId string) ([]FinalScoreBoardAdmin, error) {
	var finalScoreBoardData []FinalScoreBoardAdmin

	UserQuizResponseTable := "user_quiz_responses"
	UserPlayedQuizTable := "user_played_quizzes"

	err := model.db.
		Select(
			"users.username",
			goqu.SUM("user_quiz_responses.calculated_score").As("score"),
		).
		From("users").
		InnerJoin(goqu.T("user_played_quizzes"), goqu.On(goqu.Ex{"users.id": goqu.I("user_played_quizzes.user_id")})).
		InnerJoin(goqu.T("active_quizzes"), goqu.On(goqu.Ex{"user_played_quizzes.active_quiz_id": goqu.I("active_quizzes.id")})).
		InnerJoin(goqu.T("quizzes"), goqu.On(goqu.Ex{"active_quizzes.quiz_id": goqu.I("quizzes.id")})).
		InnerJoin(goqu.T("quiz_questions"), goqu.On(goqu.Ex{"quizzes.id": goqu.I("quiz_questions.quiz_id")})).
		InnerJoin(goqu.T("questions"), goqu.On(goqu.Ex{"quiz_questions.question_id": goqu.I("questions.id")})).
		InnerJoin(goqu.T("user_quiz_responses"), goqu.On(goqu.Ex{"questions.id": goqu.I("user_quiz_responses.question_id")})).
		Where(
			goqu.Ex{
				UserQuizResponseTable + ".user_played_quiz_id": goqu.I(UserPlayedQuizTable + ".id"),
				UserPlayedQuizTable + ".active_quiz_id":        activeQuizId,
			},
		).
		GroupBy("users.username").
		Order(goqu.I("score").Desc()).
		ScanStructs(&finalScoreBoardData)

	if err != nil {
		return nil, err
	}

	return finalScoreBoardData, nil
}
