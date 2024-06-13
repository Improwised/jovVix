package models

import (
	"database/sql"
	"encoding/json"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/doug-martin/goqu/v9"
)

type AnalyticsBoardUser struct {
	UserName         string            `db:"username" json:"username"`
	SelectedAnswer   sql.NullString    `db:"selected_answer,omitempty" json:"selected_answer"`
	CorrectAnswer    string            `db:"correct_answer,omitempty" json:"correct_answer"`
	CalculatedScore  int               `db:"calculated_score,omitempty" json:"calculated_score"`
	IsAttend         bool              `db:"is_attend,omitempty" json:"is_attend"`
	ResponseTime     int               `db:"response_time,omitempty" json:"response_time"`
	CalculatedPoints int               `db:"calculated_points,omitempty" json:"calculated_points"`
	Question         string            `db:"question,omitempty" json:"question"`
	RawOptions       []byte            `db:"options,omitempty" json:"raw_options"`
	Options          map[string]string `db:"omitempty" json:"options"`
}

type AnalyticsBoardUserModel struct {
	db *goqu.Database
}

func InitAnalyticsBoardUserModel(goqu *goqu.Database) (AnalyticsBoardUserModel, error) {
	return AnalyticsBoardUserModel{
		db: goqu,
	}, nil
}

// GetScoreForAdmin to send final score after quiz over

func (model *AnalyticsBoardUserModel) GetAnalyticsForUser(userPlayedQuizId string) ([]AnalyticsBoardUser, error) {
	var analyticsBoardData []AnalyticsBoardUser

	err := model.db.
		From(goqu.T(constants.UserQuizResponsesTable)).
		Select(
			"username",
			goqu.I(constants.UserQuizResponsesTable+".answers").As("selected_answer"),
			goqu.I(constants.QuestionsTable+".answers").As("correct_answer"),
			"calculated_score",
			"is_attend",
			"response_time",
			"calculated_points",
			"question",
			"options",
		).
		InnerJoin(goqu.T(constants.QuestionsTable), goqu.On(goqu.I(constants.UserQuizResponsesTable+".question_id").Eq(goqu.I(constants.QuestionsTable+".id")))).
		InnerJoin(goqu.T(constants.UserPlayedQuizzesTable), goqu.On(goqu.I(constants.UserPlayedQuizzesTable+".id").Eq(goqu.I(constants.UserQuizResponsesTable+".user_played_quiz_id")))).
		InnerJoin(goqu.T(constants.UsersTable), goqu.On(goqu.I(constants.UsersTable+".id").Eq(goqu.I(constants.UserPlayedQuizzesTable+".user_id")))).
		Where(goqu.Ex{
			constants.UserQuizResponsesTable + ".user_played_quiz_id": userPlayedQuizId,
		}).
		ScanStructs(&analyticsBoardData)

	if err != nil {
		return nil, err
	}
	for index := 0; index < len(analyticsBoardData); index++ {
		json.Unmarshal(analyticsBoardData[index].RawOptions, &analyticsBoardData[index].Options)
	}

	return analyticsBoardData, nil
}
