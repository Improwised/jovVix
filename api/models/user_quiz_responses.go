package models

import (
	cryptorand "crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
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
	OptionOrder      []byte    `json:"option_order,omitempty" db:"option_order"`
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
		Select(goqu.I("qq.question_id"), goqu.I("q.options"), goqu.I("q.type")).
		Join(goqu.T(QuestionTable).As("q"), goqu.On(goqu.I("q.id").Eq(goqu.I("qq.question_id")))).
		LeftJoin(goqu.T(UserQuizResponsesTable).As("uqr"), goqu.On(goqu.Ex{
			"uqr.user_played_quiz_id": userPlayedQuizId,
			"uqr.question_id":         goqu.I("qq.question_id"),
		})).
		Where(
			goqu.I("qq.active_quiz_id").Eq(quizId),
			goqu.I("uqr.id").IsNull(),
		).Executor().Query()

	if err != nil {
		return err
	}
	defer rows.Close()

	userQuizResponses := []goqu.Record{}

	for rows.Next() {
		var questionID uuid.UUID
		var rawOptions []byte
		var questionType int

		if err := rows.Scan(&questionID, &rawOptions, &questionType); err != nil {
			return err
		}

		id, err := uuid.NewUUID()

		if err != nil {
			return err
		}

		// Keep every record's keys identical. goqu builds one INSERT column list
		// for a batch, so omitting option_order for survey questions can make the
		// entire initial response copy fail.
		record := goqu.Record{
			"id":                  id,
			"question_id":         questionID,
			"user_played_quiz_id": userPlayedQuizId,
			"option_order":        nil,
		}
		if questionType == constants.SingleAnswer {
			var options map[string]string
			if err := json.Unmarshal(rawOptions, &options); err != nil {
				return err
			}

			optionOrder, err := randomOptionOrder(options)
			if err != nil {
				return err
			}
			encodedOptionOrder, err := json.Marshal(optionOrder)
			if err != nil {
				return err
			}
			record["option_order"] = string(encodedOptionOrder)
		}

		userQuizResponses = append(userQuizResponses, record)
	}
	if len(userQuizResponses) == 0 {
		return nil
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

// randomOptionOrder returns a map from the participant-visible option position to
// the original question option key. The original keys remain authoritative for scoring.
func randomOptionOrder(options map[string]string) (map[string]int, error) {
	keys := make([]int, 0, len(options))
	for key := range options {
		optionKey, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("invalid option key %q", key)
		}
		keys = append(keys, optionKey)
	}

	for index := len(keys) - 1; index > 0; index-- {
		randomIndex, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(index+1)))
		if err != nil {
			return nil, err
		}
		selectedIndex := int(randomIndex.Int64())
		keys[index], keys[selectedIndex] = keys[selectedIndex], keys[index]
	}

	optionOrder := make(map[string]int, len(keys))
	for index, originalKey := range keys {
		optionOrder[strconv.Itoa(index+1)] = originalKey
	}
	return optionOrder, nil
}

// GetOptionsForPlayer applies a saved MCQ option order. A missing mapping is a
// backward-compatible fallback for historical responses and survey questions.
func (model *UserQuizResponseModel) GetOptionsForPlayer(userPlayedQuizId, questionId uuid.UUID, options map[string]string) (map[string]string, error) {
	var rawOptionOrder sql.NullString
	found, err := model.db.From(UserQuizResponsesTable).
		Select("option_order").
		Where(goqu.Ex{"user_played_quiz_id": userPlayedQuizId, "question_id": questionId}).
		ScanVal(&rawOptionOrder)
	if err != nil {
		return nil, err
	}
	if !found || !rawOptionOrder.Valid || rawOptionOrder.String == "" {
		return options, nil
	}

	var optionOrder map[string]int
	if err := json.Unmarshal([]byte(rawOptionOrder.String), &optionOrder); err != nil {
		return nil, err
	}

	shuffledOptions := make(map[string]string, len(optionOrder))
	for displayedKey, originalKey := range optionOrder {
		option, ok := options[strconv.Itoa(originalKey)]
		if !ok {
			return nil, fmt.Errorf("option order references missing option %d", originalKey)
		}
		shuffledOptions[displayedKey] = option
	}
	return shuffledOptions, nil
}

// TranslateAnswerKeys converts a participant-visible option key into the
// original question option key before the existing score calculation runs.
func (model *UserQuizResponseModel) TranslateAnswerKeys(userPlayedQuizId, questionId uuid.UUID, answerKeys []int) ([]int, error) {
	var rawOptionOrder sql.NullString
	found, err := model.db.From(UserQuizResponsesTable).
		Select("option_order").
		Where(goqu.Ex{"user_played_quiz_id": userPlayedQuizId, "question_id": questionId}).
		ScanVal(&rawOptionOrder)
	if err != nil {
		return nil, err
	}
	if !found || !rawOptionOrder.Valid || rawOptionOrder.String == "" {
		return answerKeys, nil
	}

	var optionOrder map[string]int
	if err := json.Unmarshal([]byte(rawOptionOrder.String), &optionOrder); err != nil {
		return nil, err
	}

	translatedKeys := make([]int, len(answerKeys))
	for index, displayedKey := range answerKeys {
		originalKey, ok := optionOrder[strconv.Itoa(displayedKey)]
		if !ok {
			return nil, fmt.Errorf("invalid displayed option key %d", displayedKey)
		}
		translatedKeys[index] = originalKey
	}
	return translatedKeys, nil
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
