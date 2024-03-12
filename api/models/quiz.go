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

func (model *QuizModel) GetSharedQuestions(invitationCode int) ([]Question, sql.NullTime, error) {

	var QuestionDeliveryTime sql.NullTime = sql.NullTime{}
	statement, err := model.db.Prepare(`
	with core as (
		select
			q.*,
			aq.current_question,
			aq.is_question_active,
			aq.question_delivery_time,
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
		select
			order_no
		from
			(
			select
			(case
				when is_question_active is null then 0
				when is_question_active then order_no
				else order_no + 1
				end)
			as order_no
			from
				core
			where
				id = current_question
			union
				select
					0
		) x
		order by
			order_no desc
		limit 1
		)select
			id,
			order_no,
			question_delivery_time,
			question,
			options,
			answers,
			score,
			duration_in_seconds,
			created_at,
			updated_at
		from
			core
		where
			order_no >= (
			select
				order_no
			from
				max_order
		)
		order by
			order_no;
	`)

	if err != nil {
		return nil, QuestionDeliveryTime, err
	}

	rows, err := statement.Query(invitationCode)
	var questions []Question = []Question{}

	if err != nil {
		if err == sql.ErrNoRows {
			return questions, QuestionDeliveryTime, nil
		}

		return nil, QuestionDeliveryTime, err
	}

	for rows.Next() {
		question := Question{}
		var options []byte
		var answers []byte
		err := rows.Scan(&question.ID, &question.OrderNumber, &QuestionDeliveryTime, &question.Question, &options, &answers, &question.Score, &question.DurationInSeconds, &question.CreatedAt, &question.UpdatedAt)
		if err != nil {

			return nil, QuestionDeliveryTime, err
		}

		err = json.Unmarshal(options, &question.Options)

		if err != nil {
			return questions, QuestionDeliveryTime, err
		}

		err = json.Unmarshal(answers, &question.Answers)

		if err != nil {
			return questions, QuestionDeliveryTime, err
		}

		questions = append(questions, question)
	}

	return questions, QuestionDeliveryTime, nil
}

func (model *QuizModel) UpdateCurrentQuestion(sessionId, questionID uuid.UUID, isActive bool) error {
	records := goqu.Record{
		"current_question":   questionID,
		"is_question_active": isActive,
		"updated_at":         goqu.L("now()"),
	}

	if isActive {
		records["question_delivery_time"] = goqu.L("now()")
	} else {
		records["question_delivery_time"] = nil
	}

	result, err := model.db.Update("active_quizzes").Set(records).Where(goqu.I("id").Eq(sessionId)).Executor().Exec()

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

func (model *QuizModel) IsAllAnswerGathered(sessionId uuid.UUID) bool {
	// model.db.Select().From().Join().
	return false
}