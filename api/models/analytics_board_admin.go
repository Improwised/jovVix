package models

import (
	"database/sql"
	"encoding/json"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/doug-martin/goqu/v9"
)

type AnalyticsBoardAdmin struct {
	UserName         string            `db:"username" json:"username"`
	FirstName        string            `db:"first_name" json:"firstname"`
	SelectedAnswer   sql.NullString    `db:"selected_answer,omitempty" json:"selected_answer"`
	CorrectAnswer    string            `db:"correct_answer,omitempty" json:"correct_answer"`
	CalculatedScore  int               `db:"calculated_score,omitempty" json:"calculated_score"`
	IsAttend         bool              `db:"is_attend,omitempty" json:"is_attend"`
	ResponseTime     int               `db:"response_time,omitempty" json:"response_time"`
	CalculatedPoints int               `db:"calculated_points,omitempty" json:"calculated_points"`
	Question         string            `db:"question,omitempty" json:"question"`
	RawOptions       []byte            `db:"options,omitempty" json:"raw_options"`
	Options          map[string]string `db:"omitempty" json:"options"`
	QuestionsMedia   string            `db:"question_media" json:"question_media"`
	OptionsMedia     string            `db:"options_media" json:"options_media"`
	Resource         string            `db:"resource" json:"resource"`
	Points           int               `db:"points,omitempty" json:"points"`
	QuestionTypeID   int               `db:"type,omitempty" json:"question_type_id"`
	QuestionType     string            `db:"omitempty" json:"question_type"`
	OrderNo          int               `db:"order_no" json:"order_no"`
}

type AnalyticsBoardAdminModel struct {
	db *goqu.Database
}

func InitAnalyticsBoardAdminModel(goqu *goqu.Database) (AnalyticsBoardAdminModel, error) {
	return AnalyticsBoardAdminModel{
		db: goqu,
	}, nil
}

// GetScoreForAdmin to send final score after quiz over

func (model *AnalyticsBoardAdminModel) GetAnalyticsForAdmin(activeQuizId string) ([]AnalyticsBoardAdmin, error) {
	var analyticsBoardData []AnalyticsBoardAdmin

	err := model.db.
		From(goqu.T(constants.UserQuizResponsesTable)).
		Select(
			"username",
			"first_name",
			goqu.I(constants.UserQuizResponsesTable+".answers").As("selected_answer"),
			goqu.I(constants.QuestionsTable+".answers").As("correct_answer"),
			"calculated_score",
			"is_attend",
			"response_time",
			"calculated_points",
			"question",
			"options",
			"question_media",
			"options_media",
			"resource",
			"points",
			"type",
			"order_no",
		).
		InnerJoin(goqu.T(constants.QuestionsTable), goqu.On(goqu.I(constants.UserQuizResponsesTable+".question_id").Eq(goqu.I(constants.QuestionsTable+".id")))).
		InnerJoin(goqu.T(constants.UserPlayedQuizzesTable), goqu.On(goqu.I(constants.UserPlayedQuizzesTable+".id").Eq(goqu.I(constants.UserQuizResponsesTable+".user_played_quiz_id")))).
		InnerJoin(goqu.T(constants.UsersTable), goqu.On(goqu.I(constants.UsersTable+".id").Eq(goqu.I(constants.UserPlayedQuizzesTable+".user_id")))).
		InnerJoin(goqu.T(constants.ActiveQuizQuestionsTable), goqu.On(goqu.I(constants.UserPlayedQuizzesTable+".active_quiz_id").Eq(goqu.I(constants.ActiveQuizQuestionsTable+".active_quiz_id")), goqu.I(constants.QuestionsTable+".id").Eq(goqu.I(constants.ActiveQuizQuestionsTable+".question_id")))).
		Where(goqu.Ex{
			constants.UserPlayedQuizzesTable + ".active_quiz_id": activeQuizId,
		}).
		Order(goqu.I(constants.ActiveQuizQuestionsTable + ".order_no").Asc()).
		ScanStructs(&analyticsBoardData)

	if err != nil {
		return nil, err
	}
	for index := 0; index < len(analyticsBoardData); index++ {
		json.Unmarshal(analyticsBoardData[index].RawOptions, &analyticsBoardData[index].Options)

		analyticsBoardData[index].QuestionType, err = quizUtilsHelper.GetQuestionType(analyticsBoardData[index].QuestionTypeID)
		if err != nil {
			return nil, err
		}
	}

	return analyticsBoardData, nil
}
