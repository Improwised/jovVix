package models

import (
	"database/sql"
	"encoding/json"
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

func (model *QuizModel) GetSharedQuestions(invitationCode int) ([]Question, error) {

	statement, err := model.db.Prepare(`
	with core as (
		select
			q.*,
			aq.current_question,
			aq.is_question_active,
			aqq.order_no
		from
			active_quiz_questions aqq
		join active_quizzes aq on
			aq.invitation_code = $1
			and aq.id = aqq.active_quiz_id
		join questions q on
			q.id = aqq.question_id
		),
		max_order as (
			select order_no from (
				select order_no from core
					where id = current_question
				union
				select 0
			) x order by order_no desc limit 1
		)
		select
			id,
			order_no,
			question,
			options,
			answers,
			score,
			created_at,
			updated_at
		from
			core
		where
			order_no > (
			select
				order_no
			from
				max_order
			)
		order by
			order_no;
	`)

	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(invitationCode)
	var questions []Question = []Question{}

	if err != nil {
		if err == sql.ErrNoRows {
			return questions, nil
		}

		return nil, err
	}

	for rows.Next() {
		question := Question{}
		var options []byte
		var answers []byte
		err := rows.Scan(&question.ID, &question.OrderNumber, &question.Question, &options, &answers, &question.Score, &question.CreatedAt, &question.UpdatedAt)
		if err != nil {

			return nil, err
		}

		err = json.Unmarshal(options, &question.Options)

		if err != nil {
			return questions, err
		}

		err = json.Unmarshal(answers, &question.Answers)

		if err != nil {
			return questions, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (model *QuizModel) UpdateCurrentQuestion(sessionId, questionID uuid.UUID, isActive bool) error {
	result, err := model.db.Update("active_quizzes").Set(goqu.Record{
		"current_question":   questionID,
		"is_question_active": isActive,
	}).Where(goqu.I("id").Eq(sessionId)).Executor().Exec()

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
