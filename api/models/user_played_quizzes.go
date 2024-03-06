package models

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const UserPlayedQuizTable = "user_played_quizzes"

// QuizSession model
type UserPlayedQuiz struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	IsHost       bool      `json:"is_host" db:"is_host"`
	ActiveQuizId uuid.UUID `json:"active_quiz_id" db:"active_quiz_id"`
	LeaveAt      time.Time `json:"leave_at,omitempty" db:"leave_at"`
	QuizAnalysis string    `json:"quiz_analysis,omitempty" db:"quiz_analysis"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" db:"updated_at"`
	Status       string
}

// QuizSessionModel implements quiz session related database operations
type UserPlayedQuizModel struct {
	db          *goqu.Database
	defaultUUID uuid.UUID
}

func InitUserPlayedQuizModel(db *goqu.Database) (*UserPlayedQuizModel, error) {
	var uuid, err = uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &UserPlayedQuizModel{
		db:          db,
		defaultUUID: uuid,
	}, nil
}

// return: (uuid, int, err) -> user_quiz_session_id, status, err
// status:
//
//	0 -> error
//	1 -> query executed but parse err
//	2 -> new entry
//	3 -> existing entry
func (model *UserPlayedQuizModel) CreateUserPlayedQuizIfNotExists(userId string, quizSessionId uuid.UUID) (uuid.UUID, int, error) {

	stmt, err := model.db.Prepare(`
		with user_session_id_is_host as (
			SELECT id, is_host, 'exists' as status
			FROM user_played_quizzes us
			WHERE us.user_id = $1
			AND us.active_quiz_id = $2
			limit 1
		), insert_user_into_session as (
			INSERT INTO user_played_quizzes
				(id, user_id, active_quiz_id, is_host)
			SELECT
			$3 as id,
			$1 as user_id,
			$2 as active_quiz_id,
			(select qs.admin_id = $1 from active_quizzes qs where qs.id = $2) as is_host
			WHERE NOT EXISTS (
				select 1 from user_session_id_is_host
			) returning id, is_host, 'new' as status
		)
		select id, status from user_session_id_is_host
		union all
		select id, status from insert_user_into_session;
	`)
	if err != nil {
		return model.defaultUUID, 0, err
	}
	defer stmt.Close()

	id, err := uuid.NewUUID()

	if err != nil {
		return model.defaultUUID, 0, err
	}

	// Execute the prepared statement
	rows, err := stmt.Query(userId, quizSessionId, id)
	if err != nil {
		return model.defaultUUID, 0, err
	}
	defer rows.Close()

	var userPlayedQuizId uuid.UUID
	var status string
	if rows.Next() {
		if err := rows.Scan(&userPlayedQuizId, &status); err != nil {
			return id, 1, err
		}
	}

	if status == "new" {
		return userPlayedQuizId, 2, nil
	} else {
		return userPlayedQuizId, 3, nil
	}
}

func (model *UserPlayedQuizModel) CreateUserPlayedQuiz(userId sql.NullString, activeQuizId uuid.UUID, isHost bool) (uuid.UUID, error) {

	id, err := uuid.NewUUID()

	if err != nil {
		return model.defaultUUID, err
	}

	var userPlayedQuizId uuid.UUID
	found, err := model.db.Insert(UserPlayedQuizTable).Rows(
		goqu.Record{
			"id":             id,
			"user_id":        userId,
			"active_quiz_id": activeQuizId,
			"is_host":        isHost,
		},
	).Returning(goqu.I("id")).Executor().ScanVal(&userPlayedQuizId)

	if err != nil {
		return model.defaultUUID, err
	}

	if !found {
		return model.defaultUUID, sql.ErrNoRows
	}

	return userPlayedQuizId, nil

}
