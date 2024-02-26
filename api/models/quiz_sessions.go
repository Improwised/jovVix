package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// QuizSession model
type QuizSession struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Code             int            `json:"code" db:"code"`
	Title            string         `json:"title,omitempty" db:"title"`
	QuizID           uuid.UUID      `json:"quiz_id" db:"quiz_id"`
	AdminID          string         `json:"admin_id,omitempty" db:"admin_id"`
	ActivatedTo      sql.NullTime   `json:"activated_to,omitempty" db:"activated_to"`
	ActivatedFrom    sql.NullTime   `json:"activated_from,omitempty" db:"activated_from"`
	IsActive         bool           `json:"is_active" db:"is_active"`
	QuizAnalysis     sql.NullString `json:"quiz_analysis,omitempty" db:"quiz_analysis"`
	Current_question uuid.UUID      `json:"current_question" db:"current_question"`
	IsQuestionActive sql.NullBool   `json:"is_question_active" db:"is_question_active"`
	CreatedAt        time.Time      `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt        time.Time      `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// QuizSessionModel implements quiz session related database operations
type QuizSessionModel struct {
	db          *goqu.Database
	defaultUUID uuid.UUID
}

// InitQuizSessionModel initializes the QuizSessionModel
func InitQuizSessionModel(goqu *goqu.Database) (*QuizSessionModel, error) {
	var uuid, err = uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &QuizSessionModel{db: goqu, defaultUUID: uuid}, nil
}

func (model *QuizSessionModel) CreateQuizSession(code int, title string, quizID uuid.UUID, adminID string, maxAttempt int, activatedTo sql.NullTime, activatedFrom sql.NullTime) (uuid.UUID, error) {

	if activatedFrom.Valid && activatedFrom.Time.Before(time.Now()) {
		return model.defaultUUID, fmt.Errorf("session can not start with %s", activatedTo.Time)
	}

	if activatedFrom.Valid && activatedTo.Time.Before(activatedTo.Time) {
		return model.defaultUUID, fmt.Errorf("can not ends session before starting")
	}

	id, err := uuid.NewUUID()

	if err != nil {
		return model.defaultUUID, err
	}

	record := goqu.Record{
		"id":             id,
		"code":           code,
		"title":          title,
		"quiz_id":        quizID,
		"admin_id":       adminID,
		"max_attempt":    maxAttempt,
		"activated_to":   activatedTo,
		"activated_from": activatedFrom,
	}

	if activatedFrom.Valid {
		record["activated_from"] = nil
		record["activated_to"] = nil
	}

	if activatedTo.Valid {
		record["activated_to"] = nil
	}

	_, err = model.db.Insert("quiz_sessions").Rows(record).Executor().Exec()

	if err != nil {
		return model.defaultUUID, err
	}

	return id, nil
}

func (model *QuizSessionModel) IsUserHost(userId string, sessionId string) (bool, error) {
	query := model.db.From("quiz_sessions").
		Select(goqu.C("admin_id").Eq(userId).As("is_host")).
		Where(goqu.I("id").Eq(sessionId))

	// Execute the query
	var isHost bool
	found, err := query.ScanVal(&isHost)
	if err != nil {
		return false, err
	}

	if !found {
		return false, sql.ErrNoRows
	}

	return isHost, nil
}

type QuizSessionActive struct {
	IsActive      bool      `db:"is_active"`
	Code          int       `db:"code"`
	ActivatedFrom time.Time `db:"activated_from"`
}

func (model *QuizSessionModel) GetActiveSession(sessionId string) (QuizSession, error) {
	var quizSession QuizSession = QuizSession{}
	var isOk bool = false

	transactionObj, err := model.db.Begin()

	if err != nil {
		return quizSession, err
	}

	defer func() {
		if isOk {
			err = transactionObj.Commit()
			fmt.Println(err)
		} else {
			err = transactionObj.Rollback()
			fmt.Println(err)
		}
	}()

	found, err := model.db.Select("*").From("quiz_sessions").Where(goqu.I("id").Eq(sessionId)).Limit(1).ScanStruct(&quizSession)
	fmt.Println(isOk, sessionId, err, "--------------------")

	if err != nil {
		return quizSession, err
	}

	if !found {
		return quizSession, fmt.Errorf(constants.ErrSessionNotFound)
	}

	if quizSession.IsActive {
		return quizSession, nil
	}

	if quizSession.ActivatedTo.Valid {
		return quizSession, fmt.Errorf("session was completed")
	}

	for {
		code := quizUtilsHelper.GenerateRandomInt(100000, 999999)

		if err != nil {
			return quizSession, err
		}

		count, err := model.db.From("quiz_sessions").Where(goqu.Ex{
			"code": code,
		}).Count()

		if err != nil {
			return quizSession, err
		}

		if count == 0 {
			rows, err := model.db.Update("quiz_sessions").Set(goqu.Record{
				"code":         code,
				"is_active":    true,
				"activated_to": goqu.Literal("NOW()"),
				"updated_at":   goqu.Literal("NOW()"),
			}).Where(goqu.C("id").Eq(sessionId)).Returning("quiz_sessions.*").Executor().Query()

			if err != nil {
				return quizSession, err
			}
			defer rows.Close()

			if rows.Next() {
				err = rows.Scan(&quizSession.ID, &quizSession.Code, &quizSession.Title, &quizSession.QuizID, &quizSession.AdminID, &quizSession.ActivatedTo, &quizSession.ActivatedFrom, &quizSession.IsActive, &quizSession.QuizAnalysis, &quizSession.CreatedAt, &quizSession.UpdatedAt)
				if err != nil {
					return quizSession, err
				}
				isOk = true
				return quizSession, nil
			}
		}
	}
}
