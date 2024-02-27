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
	Code             sql.NullInt32  `json:"code" db:"code"`
	Title            string         `json:"title,omitempty" db:"title"`
	QuizID           uuid.UUID      `json:"quiz_id" db:"quiz_id"`
	AdminID          string         `json:"admin_id,omitempty" db:"admin_id"`
	ActivatedTo      sql.NullTime   `json:"activated_to,omitempty" db:"activated_to"`
	ActivatedFrom    sql.NullTime   `json:"activated_from,omitempty" db:"activated_from"`
	IsActive         bool           `json:"is_active" db:"is_active"`
	QuizAnalysis     sql.NullString `json:"quiz_analysis,omitempty" db:"quiz_analysis"`
	CurrentQuestion  sql.NullString `json:"current_question" db:"current_question"`
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

type QuizSessionActive struct {
	IsActive      bool      `db:"is_active"`
	Code          int       `db:"code"`
	ActivatedFrom time.Time `db:"activated_from"`
}

func (model *QuizSessionModel) GetSessionByCode(code string) (QuizSession, error) {
	var quizSession QuizSession = QuizSession{}

	found, err := model.db.Select("*").From("quiz_sessions").Where(goqu.I("code").Eq(code), goqu.I("is_active").Eq(true)).Limit(1).ScanStruct(&quizSession)

	if err != nil {
		return quizSession, err
	}

	if !found {
		return quizSession, fmt.Errorf(constants.ErrSessionNotFound)
	}

	return quizSession, nil
}

func (model *QuizSessionModel) GetActiveSession(sessionId string, userId string) (QuizSession, error) {
	var quizSession QuizSession = QuizSession{}
	var isOk bool = false

	transactionObj, err := model.db.Begin()

	if err != nil {
		return quizSession, err
	}

	defer func() {
		if isOk {
			err = transactionObj.Commit()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = transactionObj.Rollback()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	quizSession, err = model.GetSessionById(transactionObj, sessionId)

	if err != nil {
		return quizSession, err
	}

	if quizSession.AdminID != userId {
		return quizSession, fmt.Errorf(constants.Unauthenticated)
	}

	if quizSession.IsActive {
		return quizSession, nil
	}

	if quizSession.ActivatedTo.Valid {
		return quizSession, fmt.Errorf("session was completed")
	}

	statement, err := transactionObj.Prepare(`
	update quiz_sessions
		SET
			code=$3,
			is_active=true,
			activated_from=now(),
			updated_at = now()
		WHERE
			id=$1 and
			admin_id=$2 and
			is_active = false and
			not exists (
				select 1 from quiz_sessions where code = $3 limit 1
			)
		returning
			code
	`)

	if err != nil {
		return quizSession, err
	}

	maxTry := 10
	// handle code generation
	_, err = activateSession(maxTry, statement, quizSession.ID, userId)

	if err != nil {
		return quizSession, err
	}

	quizSession, err = model.GetSessionById(transactionObj, sessionId)
	if err != nil {
		return quizSession, err
	}
	return quizSession, nil

}

func (model *QuizSessionModel) GetSessionById(db *goqu.TxDatabase, sessionId string) (QuizSession, error) {
	var quizSession QuizSession = QuizSession{}
	found, err := db.Select("*").From("quiz_sessions").Where(goqu.I("id").Eq(sessionId)).Limit(1).ScanStruct(&quizSession)

	if err != nil {
		return quizSession, err
	}

	if !found {
		return quizSession, fmt.Errorf(constants.ErrSessionNotFound)
	}

	return quizSession, nil
}

func activateSession(maxTry int, statement *sql.Stmt, sessionId uuid.UUID, userId string) (int, error) {
	var err error
	var code int
	for {
		code = quizUtilsHelper.GenerateRandomInt(100000, 999999)

		err = statement.QueryRow(sessionId, userId, code).Scan(&code)

		if err != nil {
			if err == sql.ErrNoRows {
				maxTry -= 1
				if maxTry == 0 {
					return -1, fmt.Errorf(constants.UnknownError)
				}
				continue
			}
			return -1, err
		}

		return code, nil
	}
}
