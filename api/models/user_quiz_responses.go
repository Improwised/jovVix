package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const UserQuizResponsesTable = "user_quiz_responses"

// Question model
type UserQuizResponse struct {
	ID               uuid.UUID `json:"id" db:"id"`
	QuestionID       uuid.UUID `json:"question_id" db:"question_id"`
	Answers          []int     `json:"answers" db:"answers"`
	CalculatedScore  int       `json:"calculated_score,omitempty" db:"calculated_score"`
	IsCount          bool      `json:"is_count" db:"is_count"`
	ResponseTime     int       `json:"response_time" db:"response_time"`
	UserPlayedQuizId uuid.UUID `json:"user_played_quiz_id" db:"user_played_quiz_id"`
	CreatedAt        time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	CalculatedPoints int       `json:"calculated_points,omitempty" db:"calculated_points"`
}

type UsersQustionResponse struct {
	UserId  string         `json:"id" db:"user_id"`
	Answers sql.NullString `json:"answers" db:"answers"`
}

// QuestionModel implements question related database operations
type UserQuizResponseModel struct {
	db *goqu.Database
}

// InitQuestionModel initializes the QuestionModel
func InitUserQuizResponseModel(goqu *goqu.Database) *UserQuizResponseModel {
	return &UserQuizResponseModel{db: goqu}
}

func (model *UserQuizResponseModel) GetQuestionsCopy(userPlayedQuizId uuid.UUID, quizId uuid.UUID) error {

	rows, err := model.db.From(goqu.T("active_quiz_questions").As("qq")).
		Select(goqu.I("qq.question_id")).
		Where(goqu.I("qq.active_quiz_id").Eq(quizId)).Executor().Query()

	if err != nil {
		return err
	}
	defer rows.Close()

	userQuizResponses := []goqu.Record{}

	for rows.Next() {
		var questionID uuid.UUID

		if err := rows.Scan(&questionID); err != nil {
			return err
		}

		id, err := uuid.NewUUID()

		if err != nil {
			return err
		}

		userQuizResponses = append(userQuizResponses, goqu.Record{"id": id, "question_id": questionID, "user_played_quiz_id": userPlayedQuizId})
	}

	result, err := model.db.Insert(UserQuizResponsesTable).Rows(userQuizResponses).Executor().Exec()

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (model *UserQuizResponseModel) SubmitAnswer(userPlayedQuizId uuid.UUID, answerStruct structs.ReqAnswerSubmit, points sql.NullInt16, score, streakCount int) error {

	answerArray, err := json.Marshal(answerStruct.AnswerKeys)

	if err != nil {
		return err
	}

	result, err := model.db.Update(UserQuizResponsesTable).Set(
		goqu.Record{
			"answers":           string(answerArray),
			"calculated_points": points,
			"is_attend":         points.Valid,
			"response_time":     answerStruct.ResponseTime,
			"calculated_score":  score,
			"streak_count":      streakCount,
			"updated_at":        goqu.L("now()"),
		},
	).Where(
		goqu.I("user_played_quiz_id").Eq(userPlayedQuizId),
		goqu.I("question_id").Eq(answerStruct.QuestionId),
		goqu.I("answers").Eq(nil),
	).Executor().Exec()

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (model *UserQuizResponseModel) GetUsersResponses(sessionId uuid.UUID, questionId uuid.UUID) ([]UsersQustionResponse, error) {

	var userQuestionResponses []UsersQustionResponse
	query := model.db.From(goqu.T(constants.UserQuizResponsesTable).As("uqr")).
		Select("upq.user_id", "uqr.answers").
		Join(
			goqu.T(constants.UserPlayedQuizzesTable).As("upq"),
			goqu.On(goqu.Ex{
				"uqr.user_played_quiz_id": goqu.I("upq.id"),
			}),
		).
		Where(
			goqu.Ex{
				"uqr.question_id":    questionId,
				"upq.active_quiz_id": sessionId,
			},
		)

	err := query.ScanStructs(&userQuestionResponses)

	return userQuestionResponses, err
}
