package models

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

type FinalScoreBoardAdmin struct {
	UserName     string `db:"username" json:"username"`
	Score        int    `db:"score,omitempty" json:"score"`
	ResponseTime int    `db:"response_time,omitempty" json:"response_time"`
}

type FinalScoreBoardAdminModel struct {
	db *goqu.Database
}

func InitFinalScoreBoardAdminModel(goqu *goqu.Database) (FinalScoreBoardAdminModel, error) {
	return FinalScoreBoardAdminModel{
		db: goqu,
	}, nil
}

// GetScoreForAdmin to send final score after quiz over

func (model *FinalScoreBoardAdminModel) GetScoreForAdmin(activeQuizId string) ([]FinalScoreBoardAdmin, error) {
	var finalScoreBoardData []FinalScoreBoardAdmin

	UserQuizResponseTable := "user_quiz_responses"
	UserPlayedQuizTable := "user_played_quizzes"

	err := model.db.
		From(goqu.T("users")).
		Select(
			"users.username",
			goqu.SUM("user_quiz_responses.calculated_score").As("score"),
			goqu.Func("coalesce", goqu.SUM(goqu.Case().
				When(goqu.I("user_quiz_responses.calculated_score").Gt(0), goqu.I("user_quiz_responses.response_time")).Else(0)), 0).As("response_time"),
		).
		InnerJoin(goqu.T("user_played_quizzes"), goqu.On(goqu.I("users.id").Eq(goqu.I("user_played_quizzes.user_id")))).
		InnerJoin(goqu.T("active_quizzes"), goqu.On(goqu.I("user_played_quizzes.active_quiz_id").Eq(goqu.I("active_quizzes.id")))).
		InnerJoin(goqu.T("quizzes"), goqu.On(goqu.I("active_quizzes.quiz_id").Eq(goqu.I("quizzes.id")))).
		InnerJoin(goqu.T("quiz_questions"), goqu.On(goqu.I("quizzes.id").Eq(goqu.I("quiz_questions.quiz_id")))).
		InnerJoin(goqu.T("questions"), goqu.On(goqu.I("quiz_questions.question_id").Eq(goqu.I("questions.id")))).
		InnerJoin(goqu.T("user_quiz_responses"), goqu.On(goqu.I("questions.id").Eq(goqu.I("user_quiz_responses.question_id")))).
		Where(goqu.Ex{
			UserQuizResponseTable + ".user_played_quiz_id": goqu.I(UserPlayedQuizTable + ".id"),
			UserPlayedQuizTable + ".active_quiz_id":        activeQuizId,
		}).
		GroupBy(goqu.I("users.username")).
		Order(goqu.I("score").Desc(), goqu.I("response_time").Desc()).
		ScanStructs(&finalScoreBoardData)

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return finalScoreBoardData, nil
}
