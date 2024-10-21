package models

import (
	"database/sql"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const QuizzesTable = "quizzes"

// QuizModel implements quiz related database operations
type QuizModel struct {
	db *goqu.Database
}

// InitQuizModel initializes the QuizModel
func InitQuizModel(goquDB *goqu.Database) *QuizModel {
	return &QuizModel{db: goquDB}
}

type QuizAnalysis struct {
	ID                uuid.UUID              `json:"question_id" db:"question_id"`
	Question          string                 `json:"question" db:"question"`
	Type              int                    `json:"type" db:"type"`
	Options           map[string]string      `json:"options" db:"options"`
	QuestionsMedia    string                 `db:"question_media" json:"question_media"`
	OptionsMedia      string                 `db:"options_media" json:"options_media"`
	Resource          string                 `db:"resource" json:"resource"`
	CorrectAnswers    []int                  `json:"correct_answer" db:"answers"`
	SelectedAnswers   map[string]interface{} `json:"selected_answers" db:"selected_answers"`
	DurationInSeconds int                    `json:"duration" db:"duration_in_seconds"`
	AvgResponseTime   float32                `json:"avg_response_time" db:"avg_response_time"`
}

type QuizzesAnalysis struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	Title          string         `json:"title" db:"title"`
	Description    sql.NullString `json:"description,omitempty" db:"description"`
	ActivatedTo    sql.NullTime   `json:"activated_to,omitempty" db:"activated_to"`
	ActivatedFrom  sql.NullTime   `json:"activated_from,omitempty" db:"activated_from"`
	Questions      int            `json:"questions" db:"questions"`
	Participants   int            `json:"participants" db:"participants"`
	CorrectAnswers int            `json:"correct_answers" db:"correct_answers"`
}

type QuizWithQuestions struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	Title          string         `json:"title" db:"title" validate:"required"`
	Description    sql.NullString `json:"description,omitempty" db:"description"`
	CreatorID      string         `json:"creator_id,omitempty" db:"creator_id"`
	CreatedAt      time.Time      `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt      time.Time      `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	TotalQuestions int            `json:"total_questions" db:"total_questions"`
}

func (model *QuizModel) GetQuizzesByAdmin(creator_id string) ([]QuizWithQuestions, error) {

	questionsCountSubquery := model.db.From("quiz_questions").
		Select(goqu.COUNT("question_id")).
		Where(goqu.C("quiz_id").Eq(goqu.I("quizzes.id")))

	rows, err := model.db.From("quizzes").Select(goqu.L("*"), questionsCountSubquery.As("total_questions")).Order(goqu.I("created_at").Desc()).Where(goqu.I("creator_id").Eq(creator_id)).Executor().Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var quizzes []QuizWithQuestions = []QuizWithQuestions{}

	for rows.Next() {
		var quizWithQuestions QuizWithQuestions

		err := rows.Scan(&quizWithQuestions.ID, &quizWithQuestions.Title, &quizWithQuestions.Description, &quizWithQuestions.CreatorID, &quizWithQuestions.CreatedAt, &quizWithQuestions.UpdatedAt, &quizWithQuestions.TotalQuestions)

		if err != nil {
			return quizzes, err
		}

		decodedTitle, err := url.QueryUnescape(quizWithQuestions.Title)

		if err != nil {
			return quizzes, err
		}

		quizWithQuestions.Title = decodedTitle
		quizzes = append(quizzes, quizWithQuestions)
	}
	return quizzes, nil
}

func (model *QuizModel) GetSharedQuestions(invitationCode int) ([]Question, sql.NullTime, error) {

	var QuestionDeliveryTime sql.NullTime = sql.NullTime{}
	statement, err := model.db.Prepare(`
	with core as (
		select
			q.*,
			aq.quiz_id,
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
			quiz_id,
			order_no,
			question_delivery_time,
			question,
			options,
			answers,
			points,
			duration_in_seconds,
			question_media,
			options_media,
			resource,
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
		err := rows.Scan(&question.ID, &question.QuizId, &question.OrderNumber, &QuestionDeliveryTime, &question.Question, &options, &answers, &question.Points, &question.DurationInSeconds, &question.QuestionMedia, &question.OptionsMedia, &question.Resource, &question.CreatedAt, &question.UpdatedAt)
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

func (model *QuizModel) IsAllAnswerGathered(sessionId uuid.UUID, questionId uuid.UUID) (bool, error) {

	var skippable bool
	found, err := model.db.
		From(goqu.T(UserPlayedQuizTable).As("upq")).
		InnerJoin(goqu.T(UserQuizResponsesTable).As("upr"),
			goqu.
				On(goqu.Ex{
					"user_played_quiz_id": goqu.I("upq.id"),
					"active_quiz_id":      sessionId,
					"question_id":         questionId,
				})).Select(
		goqu.COUNT("upr.id").Eq(goqu.SUM(goqu.Case().
			When(goqu.I("upr.is_attend").Eq(true), 1).
			Else(0))).
			As("is_skippable")).ScanVal(&skippable)

	if err != nil {
		return false, err
	}

	if !found {
		return false, sql.ErrNoRows
	}

	return skippable, nil
}

func (model *QuizModel) GetQuizAnalysis(activeQuizId string) ([]QuizAnalysis, error) {
	var quizAnalysis []QuizAnalysis

	// Define the main query
	quizResponseAnalysis := model.db.From(goqu.T(UserQuizResponsesTable).As("uqr")).
		InnerJoin(goqu.T(UserPlayedQuizTable).As("upq"), goqu.On(goqu.Ex{"uqr.user_played_quiz_id": goqu.I("upq.id")})).
		InnerJoin(goqu.T(UserTable).As("u"), goqu.On(goqu.Ex{"upq.user_id": goqu.I("u.id")})).
		Select(
			goqu.C("question_id"),
			goqu.L("jsonb_object_agg(?, ?)", goqu.I("u.username"), goqu.I("uqr.answers")).As("selected_answers"),
			goqu.L("avg(?)", goqu.I("response_time")).As("avg_response_time"),
		).
		Where(goqu.Ex{"upq.active_quiz_id": activeQuizId}).
		GroupBy(goqu.C("question_id").Table("uqr"))

	// Define the final query

	query := model.db.From(goqu.T(QuestionTable).As("q")).
		With("quiz_response_analysis", quizResponseAnalysis).
		InnerJoin(goqu.T("quiz_response_analysis").As("a"), goqu.On(goqu.Ex{"q.id": goqu.I("a.question_id")})).
		InnerJoin(goqu.T(ActiveQuizQuestionsTable), goqu.On(goqu.I("q.id").Eq(goqu.I(ActiveQuizQuestionsTable+".question_id")))).
		Where(goqu.Ex{"active_quiz_questions.active_quiz_id": activeQuizId}).
		Order(goqu.I(ActiveQuizQuestionsTable+".order_no").Asc()).
		Select(
			goqu.C("question_id").Table("a"),
			goqu.C("question").Table("q"),
			goqu.C("options").Table("q"),
			goqu.C("question_media").Table("q"),
			goqu.C("options_media").Table("q"),
			goqu.C("resource").Table("q"),
			goqu.C("answers").Table("q"),
			goqu.C("selected_answers").Table("a"),
			goqu.C("duration_in_seconds").Table("q"),
			goqu.C("avg_response_time").Table("a"),
			goqu.C("type").Table("q"),
		)

	rows, err := query.Executor().Query()

	if err != nil {
		if err == sql.ErrNoRows {
			return quizAnalysis, nil
		}
		return quizAnalysis, err
	}

	for rows.Next() {
		quizAnalysisRow := QuizAnalysis{}
		var options []byte
		var answers []byte
		var selectedAnswer []byte
		err := rows.Scan(&quizAnalysisRow.ID, &quizAnalysisRow.Question, &options, &quizAnalysisRow.QuestionsMedia, &quizAnalysisRow.OptionsMedia, &quizAnalysisRow.Resource, &answers, &selectedAnswer, &quizAnalysisRow.DurationInSeconds, &quizAnalysisRow.AvgResponseTime, &quizAnalysisRow.Type)
		if err != nil {

			return nil, err
		}

		err = json.Unmarshal(options, &quizAnalysisRow.Options)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(selectedAnswer, &quizAnalysisRow.SelectedAnswers)

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(answers, &quizAnalysisRow.CorrectAnswers)

		if err != nil {
			return nil, err
		}

		quizAnalysis = append(quizAnalysis, quizAnalysisRow)
	}

	return quizAnalysis, err
}

func (model *QuizModel) ListQuizzesAnalysis(name, order, orderBy, date, userId string, page int) ([]QuizzesAnalysis, int64, error) {

	var quizzesAnalysis []QuizzesAnalysis

	quizzesAnalysisInfo := model.db.From(goqu.T(ActiveQuizzesTable).As("aq")).
		InnerJoin(goqu.T(ActiveQuizQuestionsTable).As("qq"), goqu.On(goqu.Ex{"aq.id": goqu.I("qq.active_quiz_id")})).
		InnerJoin(goqu.T(UserPlayedQuizTable).As("upq"), goqu.On(goqu.Ex{"upq.active_quiz_id": goqu.I("aq.id")})).
		InnerJoin(goqu.T(UserQuizResponsesTable).As("uqr"), goqu.On(goqu.Ex{"uqr.question_id": goqu.I("qq.question_id"), "uqr.user_played_quiz_id": goqu.I("upq.id")})).
		Where(goqu.Ex{"aq.admin_id": userId, "aq.activated_to": goqu.Op{"isNot": nil}}).
		GroupBy(
			"aq.id",
			"aq.activated_from",
			"aq.activated_to",
			"aq.quiz_id").
		Select(
			"aq.id",
			"aq.quiz_id",
			"aq.activated_from",
			"aq.activated_to",
			goqu.COUNT(goqu.DISTINCT("qq.question_id")).As("questions"),
			goqu.COUNT(goqu.DISTINCT("uqr.user_played_quiz_id")).As("participants"),
			goqu.SUM(goqu.Case().
				When(goqu.I("uqr.calculated_score").Gt(0), 1).
				When(goqu.I("uqr.calculated_score").Lte(0), 0)).
				As("correct_answers"))

	query := model.db.From(goqu.T(QuizzesTable).As("q")).
		With("quiz_analysis_info", quizzesAnalysisInfo).
		InnerJoin(goqu.T("quiz_analysis_info").As("qi"), goqu.On(goqu.Ex{"qi.quiz_id": goqu.I("q.id")})).
		Select(
			"qi.questions",
			"qi.participants",
			"qi.correct_answers",
			"qi.id",
			"qi.activated_from",
			"qi.activated_to",
			"q.title",
			"q.description")

	if name != "" {
		query = query.Where(goqu.Func("LOWER", goqu.I("title")).Like("%" + strings.ToLower(name) + "%"))
	}

	if date != "" {
		query = query.Where(goqu.C("activated_from").Gte(date))
	}

	count, err := query.Count()

	if err != nil {
		return quizzesAnalysis, 0, err
	}

	if order == "desc" {
		query = query.Order(goqu.I(orderBy).Desc())
	} else {
		query = query.Order(goqu.I(orderBy).Asc())
	}

	offset := (page - 1) * constants.DefaultPageSize
	if offset >= 0 {
		query = query.Offset(uint(offset))
	}

	query = query.Limit(uint(constants.DefaultPageSize))

	sql, args, err := query.ToSQL()
	if err != nil {
		return quizzesAnalysis, 0, err
	}

	err = model.db.ScanStructs(&quizzesAnalysis, sql, args...)

	return quizzesAnalysis, count, err
}

// Delete quiz and and their question using `quiz_id`
func (model *QuizModel) DeleteQuizById(transaction *goqu.TxDatabase, QuizId string) error {
	questionIds := []string{}
	err := transaction.Delete(constants.QuizQuestionsTable).Where(goqu.Ex{"quiz_id": QuizId}).Returning("question_id").Executor().ScanVals(&questionIds)
	if err != nil {
		return err
	}

	err = deleteQuestionsByIds(transaction, questionIds)
	if err != nil {
		return err
	}

	err = deleteQuizById(transaction, QuizId)

	return err
}

// Delete Quiz by Id (only if no active quiz is present)
func deleteQuizById(transaction *goqu.TxDatabase, quizId string) error {

	_, err := transaction.Delete(QuizzesTable).Where(goqu.Ex{"id": quizId}).Executor().Exec()
	if err != nil {
		return err
	}

	return nil
}

// Deletes all quizzes created by a user and their related questions.
// It deletes from the QuizQuestionsTable, removes related questions, and finally deletes the quizzes themselves.
func (model *QuizModel) DeleteCreatedQuizzesByUserId(transaction *goqu.TxDatabase, userId string) error {

	quizSubquery := transaction.From(QuizzesTable).Select("id").Where(goqu.Ex{"creator_id": userId})

	questionIds := []string{}
	err := transaction.Delete(constants.QuizQuestionsTable).Where(goqu.Ex{"quiz_id": goqu.Op{"in": quizSubquery}}).Returning("question_id").Executor().ScanVals(&questionIds)
	if err != nil {
		return err
	}

	err = deleteQuestionsByIds(transaction, questionIds)
	if err != nil {
		return err
	}

	_, err = transaction.Delete(QuizzesTable).Where(goqu.Ex{"id": goqu.Op{"in": quizSubquery}}).Executor().Exec()

	return err
}
