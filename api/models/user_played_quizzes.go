package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/pkg/structs"
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

func InitUserPlayedQuizModel(db *goqu.Database) *UserPlayedQuizModel {
	uuid := uuid.UUID{}

	return &UserPlayedQuizModel{
		db:          db,
		defaultUUID: uuid,
	}
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

func (model *UserPlayedQuizModel) GetActiveSession(id string, invitationCode string, userID string) (ActiveQuiz, error) {
	var activeQuiz ActiveQuiz
	userIDObj := sql.NullString{
		Valid:  userID != "",
		String: userID,
	}

	found, err := model.db.From(goqu.T(UserPlayedQuizTable).As("upq")).Join(
		goqu.I(ActiveQuizzesTable).As("aq"),
		goqu.On(goqu.I("upq.active_quiz_id").Eq(goqu.I("aq.id"))),
	).Select(
		goqu.I("aq.*"),
	).Where(
		goqu.Ex{
			"upq.id":             id,
			"upq.user_id":        userIDObj,
			"aq.invitation_code": invitationCode,
		},
	).Limit(1).ScanStruct(&activeQuiz)

	if err != nil {
		return activeQuiz, err
	}
	if !found {
		return activeQuiz, sql.ErrNoRows
	}

	if !activeQuiz.IsActive {
		return activeQuiz, fmt.Errorf(constants.ErrInvitationCodeNotFound)
	}

	return activeQuiz, nil
}

func (model *UserPlayedQuizModel) GetCurrentActiveQuestion(id uuid.UUID) (uuid.UUID, error) {
	var currentQuestion uuid.UUID
	found, err := model.db.Select("current_question").From(goqu.T(ActiveQuizzesTable).As("aq")).Join(goqu.T(UserPlayedQuizTable).As("upq"), goqu.On(goqu.I("upq.id").Eq(id), goqu.I("aq.is_question_active").Eq(true), goqu.I("upq.active_quiz_id").Eq(goqu.I("aq.id")))).ScanVal(&currentQuestion)

	if err != nil {
		return uuid.UUID{}, err
	}

	if !found {
		return uuid.UUID{}, sql.ErrNoRows
	}

	return currentQuestion, nil
}

type UserRank struct {
	Rank         int    `json:"rank" db:"rank"`
	Points       int    `json:"points" db:"points"`
	Score        int    `json:"score" db:"calculated_score"`
	ResponseTime int    `json:"response_time" db:"response_time"`
	UserName     string `json:"username" db:"username"`
	FirstName    string `json:"firstname" db:"first_name"`
}

func (model *UserPlayedQuizModel) GetRank(sessionId uuid.UUID, questionId uuid.UUID) ([]UserRank, error) {

	mainQuery := model.db.
		From(goqu.T("get_question_info").As("gqi")).
		Join(goqu.T("get_sum").As("gs"), goqu.On(goqu.I("gs.user_id").Eq(goqu.I("gqi.user_id")))).
		Join(goqu.T(UserTable).As("u"), goqu.On(goqu.I("gs.user_id").Eq(goqu.I("u.id"))))

	// Define the common table expressions (CTEs)
	core := mainQuery.
		With("core", goqu.
			Select("uqr.calculated_score", "uqr.calculated_points", "uqr.question_id", "uqr.response_time", "uqr.is_attend", "upq.user_id").
			From(goqu.T(UserPlayedQuizTable).As("upq")).
			Join(goqu.T("user_quiz_responses").As("uqr"), goqu.On(goqu.Ex{
				"upq.id":             goqu.I("uqr.user_played_quiz_id"),
				"upq.active_quiz_id": sessionId,
			}),
			))

	getSum := core.
		With("get_sum", goqu.
			Select("user_id", goqu.SUM("calculated_score").As("calculated_total_score"), goqu.SUM("calculated_points").As("total_points")).
			From("core").
			GroupBy("user_id"),
		)
	getQuestionInfo := getSum.
		With("get_question_info", goqu.
			Select("is_attend", "user_id", "response_time").
			From("core").
			Where(goqu.Ex{
				"question_id": questionId,
			}),
		)
	final_query := getQuestionInfo.Select(
		goqu.DENSE_RANK().Over(goqu.W().OrderBy(goqu.I("gs.calculated_total_score").Desc())).As("rank"),
		goqu.I("gs.calculated_total_score"),
		goqu.I("gs.total_points"),
		goqu.I("gqi.response_time"),
		goqu.I("u.username"),
		goqu.I("u.first_name"),
	)

	// Define the main query to filter by rank
	rows, err := final_query.Executor().Query()

	if err != nil {
		return []UserRank{}, err
	}
	defer rows.Close()

	userRanks := []UserRank{}
	for rows.Next() {
		var userRank UserRank
		err := rows.Scan(&userRank.Rank, &userRank.Score, &userRank.Points, &userRank.ResponseTime, &userRank.UserName, &userRank.FirstName)
		if err != nil {
			return userRanks, err
		}
		userRanks = append(userRanks, userRank)

	}
	return userRanks, nil
}

func (model *UserPlayedQuizModel) ListUserPlayedQuizes(userId string) ([]structs.ResUserPlayedQuiz, error) {
	var userPlayedQuiz []structs.ResUserPlayedQuiz

	query := model.db.From(UserPlayedQuizTable).
		Select(
			"quizzes.title",
			"quizzes.description",
			"user_played_quizzes.id",
			"user_played_quizzes.created_at",
			goqu.COUNT(goqu.I("quiz_questions.id")).As("total_questions"),
		).
		InnerJoin(goqu.T(constants.ActiveQuizzesTable), goqu.On(goqu.I(UserPlayedQuizTable+".active_quiz_id").Eq(goqu.I(constants.ActiveQuizzesTable+".id")))).
		InnerJoin(goqu.T(constants.QuizzesTable), goqu.On(goqu.I(ActiveQuizzesTable+".quiz_id").Eq(goqu.I(constants.QuizzesTable+".id")))).
		InnerJoin(goqu.T(constants.QuizQuestionsTable), goqu.On(goqu.I(constants.QuizzesTable+".id").Eq(goqu.I(constants.QuizQuestionsTable+".quiz_id")))).
		Where(goqu.Ex{
			UserPlayedQuizTable + ".user_id": userId,
		}).GroupBy("user_played_quizzes.id", "quizzes.id").Order(goqu.I("user_played_quizzes.created_at").Desc())

	query = query.Limit(constants.PageSize)

	sql, args, err := query.ToSQL()
	if err != nil {
		return userPlayedQuiz, err
	}

	err = model.db.ScanStructs(&userPlayedQuiz, sql, args...)

	return userPlayedQuiz, err
}

func (model *UserPlayedQuizModel) ListUserPlayedQuizesWithQuestionById(UserPlayedQuizId string) ([]structs.ResUserPlayedQuizAnalyticsBoard, error) {
	var userPlayedQuizAnalyticsBoard []structs.ResUserPlayedQuizAnalyticsBoard

	query := model.db.From(UserPlayedQuizTable).
		Select(
			goqu.I(constants.UserQuizResponsesTable+".answers").As("selected_answer"),
			goqu.I(constants.QuestionsTable+".answers").As("correct_answer"),
			"calculated_score",
			"is_attend",
			"response_time",
			"calculated_points",
			"question",
			"options",
			"points",
			"type",
		).
		InnerJoin(goqu.T(constants.UserQuizResponsesTable), goqu.On(goqu.I(UserPlayedQuizTable+".id").Eq(goqu.I(constants.UserQuizResponsesTable+".user_played_quiz_id")))).
		InnerJoin(goqu.T(constants.QuestionsTable), goqu.On(goqu.I(constants.UserQuizResponsesTable+".question_id").Eq(goqu.I(constants.QuestionsTable+".id")))).
		Where(goqu.Ex{
			UserPlayedQuizTable + ".id": UserPlayedQuizId,
		})

	sql, args, err := query.ToSQL()
	if err != nil {
		return userPlayedQuizAnalyticsBoard, err
	}

	err = model.db.ScanStructs(&userPlayedQuizAnalyticsBoard, sql, args...)
	if err != nil {
		return userPlayedQuizAnalyticsBoard, err
	}

	for index := 0; index < len(userPlayedQuizAnalyticsBoard); index++ {
		json.Unmarshal(userPlayedQuizAnalyticsBoard[index].RawOptions, &userPlayedQuizAnalyticsBoard[index].Options)

		userPlayedQuizAnalyticsBoard[index].QuestionType, err = quizUtilsHelper.GetQuestionType(userPlayedQuizAnalyticsBoard[index].QuestionTypeID)
		if err != nil {
			return nil, err
		}
	}

	return userPlayedQuizAnalyticsBoard, nil
}
