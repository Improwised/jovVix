package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type FinalScoreBoard struct {
	UserName     string `db:"username" json:"username"`
	Score        int    `db:"score" json:"score"`
	ResponseTime int    `db:"response_time" json:"response_time"`
}

type FinalScoreBoardModel struct {
	db *goqu.Database
}

func InitFinalScoreBoardModel(goqu *goqu.Database) (FinalScoreBoardModel, error) {
	return FinalScoreBoardModel{
		db: goqu,
	}, nil
}

// GetScore to send final score after quiz over

func (model *FinalScoreBoardModel) GetScore(user_played_quiz string) ([]FinalScoreBoard, error) {
	var finalScoreBoardData []FinalScoreBoard
	var activeQuizId uuid.UUID

	UserQuizResponseTable := "user_quiz_responses"
	UserPlayedQuizTable := "user_played_quizzes"

	_, err := model.db.Select("active_quiz_id").From("user_played_quizzes").Where(goqu.I("id").Eq(uuid.MustParse(user_played_quiz))).ScanVal(&activeQuizId)

	if err != nil {
		return nil, err
	}

	err = model.db.
		From(goqu.T("users")).
		Select(
			"users.username",
			goqu.SUM("user_quiz_responses.calculated_score").As("score"),
			goqu.Func("coalesce", goqu.SUM(goqu.Case().
				When(goqu.I("user_quiz_responses.calculated_score").Gt(0), goqu.I("user_quiz_responses.response_time")).Else(0)), 0).As("response_time"),
		).
		// From("users").
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
		Order(goqu.I("score").Desc(), goqu.I("response_time").Asc()).
		ScanStructs(&finalScoreBoardData)

	if err != nil {
		return nil, err
	}

	return finalScoreBoardData, nil
}
