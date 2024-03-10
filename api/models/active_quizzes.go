package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const ActiveQuizzesTable = "active_quizzes"

// ActiveQuiz model
type ActiveQuiz struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	InvitationCode   sql.NullInt32  `json:"invitation_code" db:"invitation_code"`
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

// ActiveQuizModel implements quiz session related database operations
type ActiveQuizModel struct {
	db          *goqu.Database
	defaultUUID uuid.UUID
	logger      *zap.Logger
}

// InitActiveQuizModel initializes the ActiveQuizModel
func InitActiveQuizModel(goqu *goqu.Database, logger *zap.Logger) *ActiveQuizModel {
	var uuid = uuid.UUID{}
	return &ActiveQuizModel{db: goqu, defaultUUID: uuid, logger: logger}
}

func (model *ActiveQuizModel) CreateActiveQuiz(invitationCode int, title string, quizID uuid.UUID, adminID string, maxAttempt int, activatedTo sql.NullTime, activatedFrom sql.NullTime) (uuid.UUID, error) {

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
		"id":              id,
		"invitation_code": invitationCode,
		"title":           title,
		"quiz_id":         quizID,
		"admin_id":        adminID,
		"max_attempt":     maxAttempt,
		"activated_to":    activatedTo,
		"activated_from":  activatedFrom,
	}

	if activatedFrom.Valid {
		record["activated_from"] = nil
		record["activated_to"] = nil
	}

	if activatedTo.Valid {
		record["activated_to"] = nil
	}

	_, err = model.db.Insert(ActiveQuizzesTable).Rows(record).Executor().Exec()

	if err != nil {
		return model.defaultUUID, err
	}

	return id, nil
}

func (model *ActiveQuizModel) GetSessionByCode(invitationCode string) (ActiveQuiz, error) {
	var activeQuiz ActiveQuiz = ActiveQuiz{}

	found, err := model.db.Select("*").From(ActiveQuizzesTable).Where(goqu.I("invitation_code").Eq(invitationCode), goqu.I("is_active").Eq(true)).Limit(1).ScanStruct(&activeQuiz)

	if err != nil {
		return activeQuiz, err
	}

	if !found {
		return activeQuiz, sql.ErrNoRows
	}

	return activeQuiz, nil
}

func (model *ActiveQuizModel) GetOrActivateSession(sessionId string, userId string) (ActiveQuiz, error) {
	var activeQuiz ActiveQuiz = ActiveQuiz{}
	var isOk bool = false

	transactionObj, err := model.db.Begin()

	if err != nil {
		return activeQuiz, err
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

	activeQuiz, err = model.GetSessionById(transactionObj, sessionId)

	if err != nil {
		return activeQuiz, err
	}

	if activeQuiz.AdminID != userId {
		return activeQuiz, fmt.Errorf(constants.Unauthenticated)
	}

	if activeQuiz.IsActive {
		isOk = true
		return activeQuiz, nil
	}

	if activeQuiz.ActivatedTo.Valid {
		return activeQuiz, fmt.Errorf(constants.ErrSessionWasCompleted)
	}

	maxTry := 10
	// handle invitation_code generation
	invitation_code, err := activateSession(transactionObj, maxTry, activeQuiz.ID, userId)

	if err != nil {
		return activeQuiz, err
	}
	isOk = (invitation_code != -1)

	activeQuiz, err = model.GetSessionById(transactionObj, sessionId)
	if err != nil {
		return activeQuiz, err
	}
	return activeQuiz, nil

}

func (model *ActiveQuizModel) GetSessionById(db *goqu.TxDatabase, sessionId string) (ActiveQuiz, error) {
	var activeQuiz ActiveQuiz = ActiveQuiz{}
	found, err := db.Select("*").From(ActiveQuizzesTable).Where(goqu.I("id").Eq(sessionId)).Limit(1).ScanStruct(&activeQuiz)

	if err != nil {
		return activeQuiz, err
	}

	if !found {
		return activeQuiz, fmt.Errorf(constants.ErrSessionNotFound)
	}

	return activeQuiz, nil
}

func (model *ActiveQuizModel) GetSession(sessionId string) (ActiveQuiz, error) {
	var activeQuiz ActiveQuiz = ActiveQuiz{}
	found, err := model.db.Select("*").From(ActiveQuizzesTable).Where(goqu.I("id").Eq(sessionId)).Limit(1).ScanStruct(&activeQuiz)

	if err != nil {
		return activeQuiz, err
	}

	if !found {
		return activeQuiz, fmt.Errorf(constants.ErrSessionNotFound)
	}

	return activeQuiz, nil
}

func activateSession(transactionObj *goqu.TxDatabase, maxTry int, sessionId uuid.UUID, userId string) (int, error) {
	var err error
	var invitation_code int
	statement, err := transactionObj.Prepare(`
	update active_quizzes
		SET
			invitation_code=$3,
			is_active=true,
			activated_from=now(),
			updated_at = now()
		WHERE
			id=$1 and
			admin_id=$2 and
			is_active = false and
			not exists (
				select 1 from active_quizzes where invitation_code = $3 limit 1
			)
		returning
			invitation_code
	`)

	if err != nil {
		return -1, err
	}

	defer statement.Close()

	for {
		invitation_code = quizUtilsHelper.GenerateRandomInt(constants.MinInvitationCode, constants.MaxInvitationCode)

		err = statement.QueryRow(sessionId, userId, invitation_code).Scan(&invitation_code)

		if err != nil {
			if err == sql.ErrNoRows {
				maxTry -= 1
				if maxTry == 0 {
					return -1, fmt.Errorf(constants.ErrMaxTryToGenerateCode)
				}
				continue
			}
			return -1, err
		}

		return invitation_code, nil
	}
}

func (model *ActiveQuizModel) Deactivate(id uuid.UUID) error {
	result, err := model.db.Update("active_quizzes").Set(goqu.Record{
		"invitation_code":    nil,
		"is_active":          false,
		"activated_to":       goqu.L("now()"),
		"current_question":   nil,
		"is_question_active": nil,
	}).Where(goqu.I("id").Eq(id)).Executor().Exec()

	if err != nil {
		return err
	}

	affectedRow, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (model *ActiveQuizModel) GetCurrentActiveQuestion(id uuid.UUID) (uuid.UUID, error) {
	var currentQuestion uuid.UUID
	found, err := model.db.Select("current_question").From(ActiveQuizzesTable).Where(goqu.I("id").Eq(id), goqu.I("is_question_active").Eq(true)).ScanVal(&currentQuestion)

	if err != nil {
		return uuid.UUID{}, err
	}

	if !found {
		return uuid.UUID{}, sql.ErrNoRows
	}

	return currentQuestion, nil
}
