package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

// Quiz model
type Quiz struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" validate:"required"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatorID   string    `json:"creator_id,omitempty" db:"creator_id"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// QuizModel implements quiz related database operations
type QuizModel struct {
	db *goqu.Database
}

// InitQuizModel initializes the QuizModel
func InitQuizModel(goquDB *goqu.Database) *QuizModel {
	return &QuizModel{db: goquDB}
}

type QuizActivity struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Title        string    `json:"title" db:"title" validate:"required"`
	Description  string    `json:"description,omitempty" db:"description"`
	UserActivity string    `json:"user_activity" db:"role"`
}

func (model *QuizModel) GetAllQuizzesActivity(user_id string) ([]QuizActivity, error) {
	var quizzes []QuizActivity = []QuizActivity{}

	// user as a host
	hostQuizzes := model.db.From(goqu.T("active_quizzes").As("qs")).
		Select(
			goqu.C("quiz_id"),
			goqu.L("'host'").As("user_activity"),
		).
		Where(goqu.I("qs.admin_id").Eq(user_id))

	// user as a creator
	creatorQuizzes := hostQuizzes.Union(
		model.db.From(goqu.T("quizzes").As("q")).
			Select(
				goqu.I("id"),
				goqu.L("'creator'").As("user_activity"),
			).
			Where(goqu.I("q.creator_id").Eq(user_id)),
	)

	err := model.db.Select("*").From(goqu.T("quizzes").As("q")).Join(
		creatorQuizzes.As("quiz_activity"),
		goqu.On(goqu.I("quiz_activity.quiz_id").Eq(goqu.I("q.id"))),
	).
		Select(
			goqu.I("q.id"),
			goqu.I("q.title"),
			goqu.I("q.description"),
			goqu.I("quiz_activity.user_activity"),
		).Executor().ScanStructs(&quizzes)

	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (model *QuizModel) GetQuizzesByAdmin(creator_id string) ([]Quiz, error) {
	var quizzes []Quiz = []Quiz{}

	rows, err := model.db.Select("*").From("quizzes").Where(goqu.Ex{"creator_id": creator_id}).Executor().Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var quiz Quiz
		err := rows.Scan(quiz)

		if err != nil {
			return nil, err
		}

		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

func (model *QuizModel) GetSharedQuestions(invitationCode int) ([]SessionQuestion, error) {

	statement, err := model.db.Prepare(`
	with get_ids as (
		select id, quiz_id from active_quizzes qs where qs.invitation_code = $1 and qs.is_active = $2
	)
	, get_question_order as(
		select order_no, next_question, qs.is_question_active from session_questions sq2 join get_ids ids on ids.id = sq2.active_quiz_id join active_quizzes qs on qs.current_question = sq2.question_id
	), get_questions as (
		select sq.* from session_questions sq join get_ids ids on ids.id = sq.active_quiz_id and order_no >= (
			select (
			case when is_question_active then order_no + 1
			else order_no end
			) from get_question_order
		)
	) select id, question_id, next_question, active_quiz_id, order_no from get_questions;
	`)

	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(invitationCode, true)
	var questions []SessionQuestion = []SessionQuestion{}

	if err != nil {
		if err == sql.ErrNoRows {
			return questions, nil
		}

		return nil, err
	}

	for rows.Next() {
		question := SessionQuestion{}
		err := rows.Scan(&question.ID, &question.QuestionID, &question.NextQuestion, &question.QuizSessionID, &question.OrderNo)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			return nil, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func IsValidCode(model *QuizModel) {

}
