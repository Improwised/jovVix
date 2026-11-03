package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/pkg/structs"
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

func (model *UserPlayedQuizModel) GetUserPlayedQuizID(userId string, activeQuizId uuid.UUID) (uuid.UUID, error) {
	var userPlayedQuizId uuid.UUID
	found, err := model.db.From(UserPlayedQuizTable).
		Select("id").
		Where(goqu.Ex{"user_id": userId, "active_quiz_id": activeQuizId}).
		ScanVal(&userPlayedQuizId)
	if err != nil {
		return uuid.UUID{}, err
	}
	if !found {
		return uuid.UUID{}, sql.ErrNoRows
	}
	return userPlayedQuizId, nil
}

// return: (uuid, int, err) -> user_quiz_session_id, status, err
// status:
//
//	0 -> error
//	1 -> query executed but parse err
//	2 -> new entry
//	3 -> existing entry
var ErrPlayerNameTaken = errors.New("player name is already taken")
var ErrInvalidPlayerName = errors.New("player name must be between 1 and 12 characters")

// JoinQuiz serializes joins for a session, including checks against old records
// that predate player_name. The unique index also protects new records.
func (model *UserPlayedQuizModel) JoinQuiz(userID string, sessionID uuid.UUID, requestedName string) (uuid.UUID, int, string, []string, error) {
	tx, err := model.db.Begin()
	if err != nil {
		return uuid.Nil, 0, "", nil, err
	}
	defer tx.Rollback()
	var locked uuid.UUID
	if err = tx.QueryRow("SELECT id FROM active_quizzes WHERE id = $1 FOR UPDATE", sessionID).Scan(&locked); err != nil {
		return uuid.Nil, 0, "", nil, err
	}
	var id uuid.UUID
	var savedName string
	err = tx.QueryRow(`SELECT upq.id, COALESCE(upq.player_name, u.first_name) FROM user_played_quizzes upq JOIN users u ON u.id = upq.user_id WHERE upq.user_id = $1 AND upq.active_quiz_id = $2 LIMIT 1`, userID, sessionID).Scan(&id, &savedName)
	if err == nil {
		return id, 3, savedName, nil, nil
	}
	if err != sql.ErrNoRows {
		return uuid.Nil, 0, "", nil, err
	}
	var first, last string
	if err = tx.QueryRow("SELECT first_name, COALESCE(last_name, '') FROM users WHERE id = $1", userID).Scan(&first, &last); err != nil {
		return uuid.Nil, 0, "", nil, err
	}
	name := strings.TrimSpace(requestedName)
	if name == "" {
		name = strings.TrimSpace(first)
		if chars := []rune(name); len(chars) > 12 {
			name = string(chars[:12])
		}
	}
	if utf8.RuneCountInString(name) < 1 || utf8.RuneCountInString(name) > 12 {
		return uuid.Nil, 0, "", nil, ErrInvalidPlayerName
	}
	available := func(candidate string) (bool, error) {
		var taken bool
		err := tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM user_played_quizzes upq JOIN users u ON u.id = upq.user_id WHERE upq.active_quiz_id = $1 AND upq.is_host = false AND lower(btrim(COALESCE(upq.player_name, u.first_name))) = lower(btrim($2)))`, sessionID, candidate).Scan(&taken)
		return !taken, err
	}
	ok, err := available(name)
	if err != nil {
		return uuid.Nil, 0, "", nil, err
	}
	if !ok {
		suggestions, err := suggestPlayerNames(first, last, name, available)
		if err != nil {
			return uuid.Nil, 0, "", nil, err
		}
		return uuid.Nil, 0, "", suggestions, ErrPlayerNameTaken
	}
	id, err = uuid.NewRandom()
	if err != nil {
		return uuid.Nil, 0, "", nil, err
	}
	_, err = tx.Exec(`INSERT INTO user_played_quizzes (id, user_id, active_quiz_id, is_host, player_name) VALUES ($1, $2, $3, false, $4)`, id, userID, sessionID, name)
	if err != nil {
		return uuid.Nil, 0, "", nil, err
	}
	if err = tx.Commit(); err != nil {
		return uuid.Nil, 0, "", nil, err
	}
	return id, 2, name, nil, nil
}

func suggestPlayerNames(first, last, name string, available func(string) (bool, error)) ([]string, error) {
	suggestions := []string{}
	seen := map[string]bool{strings.ToLower(name): true}
	try := func(candidate string) error {
		candidate = strings.TrimSpace(candidate)
		key := strings.ToLower(candidate)
		if candidate == "" || utf8.RuneCountInString(candidate) > 12 || seen[key] {
			return nil
		}
		seen[key] = true
		free, err := available(candidate)
		if err == nil && free {
			suggestions = append(suggestions, candidate)
		}
		return err
	}
	if strings.TrimSpace(last) != "" {
		for _, candidate := range []string{last, strings.ReplaceAll(first+last, " ", "")} {
			if err := try(candidate); err != nil {
				return nil, err
			}
		}
	}
	base := []rune(name)
	if len(base) > 9 {
		base = base[:9]
	}
	for i := 0; i < 30 && len(suggestions) < 4; i++ {
		if err := try(fmt.Sprintf("%s%d", string(base), quizUtilsHelper.GenerateRandomInt(10, 999))); err != nil {
			return nil, err
		}
	}
	return suggestions, nil
}

func (model *UserPlayedQuizModel) GetPlayerName(userID string, sessionID uuid.UUID) (string, error) {
	var name string
	found, err := model.db.From(goqu.T(UserPlayedQuizTable).As("upq")).
		Join(goqu.T(UserTable).As("u"), goqu.On(goqu.I("u.id").Eq(goqu.I("upq.user_id")))).
		Select(goqu.L("COALESCE(upq.player_name, u.first_name)")).
		Where(goqu.Ex{"upq.user_id": userID, "upq.active_quiz_id": sessionID}).ScanVal(&name)
	if err != nil {
		return "", err
	}
	if !found {
		return "", sql.ErrNoRows
	}
	return name, nil
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

func (model *UserPlayedQuizModel) GetCurrentActiveQuestion(id string) (uuid.UUID, error) {
	var currentQuestion uuid.UUID
	found, err := model.db.Select("current_question").From(ActiveQuizzesTable).Where(goqu.Ex{"is_question_active": true, "id": id}).ScanVal(&currentQuestion)

	if err != nil {
		return uuid.UUID{}, err
	}

	if !found {
		return uuid.UUID{}, sql.ErrNoRows
	}

	return currentQuestion, nil
}

// IsParticipant verifies that a played-quiz record belongs to the requesting user,
// belongs to the requested session, and is not a host record.
func (model *UserPlayedQuizModel) IsParticipant(userID string, userPlayedQuizID uuid.UUID, activeQuizID string) (bool, error) {
	count, err := model.db.From(UserPlayedQuizTable).Where(goqu.Ex{
		"id":             userPlayedQuizID,
		"user_id":        userID,
		"active_quiz_id": activeQuizID,
		"is_host":        false,
	}).Count()
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

type UserRank struct {
	Rank         int    `json:"rank" db:"rank"`
	Points       int    `json:"points" db:"points"`
	Score        int    `json:"score" db:"calculated_score"`
	ResponseTime int    `json:"response_time" db:"response_time"`
	UserName     string `json:"username" db:"username"`
	FirstName    string `json:"firstname" db:"first_name"`
	ImageKey     string `json:"img_key" db:"img_key"`
	StreakCount  int    `json:"streak_count" db:"streak_count"`
}

func (model *UserPlayedQuizModel) GetRank(sessionId uuid.UUID, questionId uuid.UUID) ([]UserRank, error) {

	mainQuery := model.db.
		From(goqu.T("get_question_info").As("gqi")).
		Join(goqu.T("get_sum").As("gs"), goqu.On(goqu.I("gs.user_id").Eq(goqu.I("gqi.user_id")))).
		Join(goqu.T(UserTable).As("u"), goqu.On(goqu.I("gs.user_id").Eq(goqu.I("u.id")))).
		Join(goqu.T(UserPlayedQuizTable).As("participant"), goqu.On(goqu.I("participant.id").Eq(goqu.I("gqi.participant_id"))))

	// Define the common table expressions (CTEs)
	core := mainQuery.
		With("core", goqu.
			Select("uqr.calculated_score", "uqr.calculated_points", "uqr.question_id", "uqr.response_time", "uqr.is_attend", "upq.user_id", goqu.I("upq.id").As("participant_id"), "uqr.streak_count").
			From(goqu.T(UserPlayedQuizTable).As("upq")).
			Join(goqu.T("user_quiz_responses").As("uqr"), goqu.On(goqu.Ex{
				"upq.id":             goqu.I("uqr.user_played_quiz_id"),
				"upq.active_quiz_id": sessionId,
			}),
			).Where(goqu.I("upq.is_host").Eq(false)))

	getSum := core.
		With("get_sum", goqu.
			Select("user_id", goqu.SUM("calculated_score").As("calculated_total_score"), goqu.SUM("calculated_points").As("total_points")).
			From("core").
			GroupBy("user_id"),
		)
	getQuestionInfo := getSum.
		With("get_question_info", goqu.
			Select("is_attend", "user_id", "participant_id", "response_time", "streak_count").
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
		goqu.L("COALESCE(participant.player_name, u.first_name)").As("first_name"),
		goqu.I("u.img_key"),
		goqu.I("gqi.streak_count"),
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
		err := rows.Scan(&userRank.Rank, &userRank.Score, &userRank.Points, &userRank.ResponseTime, &userRank.UserName, &userRank.FirstName, &userRank.ImageKey, &userRank.StreakCount)
		if err != nil {
			return userRanks, err
		}
		userRanks = append(userRanks, userRank)

	}
	return userRanks, nil
}

func (model *UserPlayedQuizModel) ListUserPlayedQuizes(userId string, page int, titleSearch string) ([]structs.ResUserPlayedQuiz, int64, error) {
	var userPlayedQuiz []structs.ResUserPlayedQuiz

	query := model.db.From(UserPlayedQuizTable).
		Select(
			"quizzes.title",
			"quizzes.description",
			"user_played_quizzes.id",
			"user_played_quizzes.created_at",
			goqu.COUNT(goqu.I("quiz_questions.id")).As("total_questions"),
			goqu.COUNT(goqu.I("user_played_quizzes.id")).Over(goqu.W().PartitionBy()).As("total_count"),
		).
		InnerJoin(goqu.T(constants.ActiveQuizzesTable), goqu.On(goqu.I(UserPlayedQuizTable+".active_quiz_id").Eq(goqu.I(constants.ActiveQuizzesTable+".id")))).
		InnerJoin(goqu.T(constants.QuizzesTable), goqu.On(goqu.I(ActiveQuizzesTable+".quiz_id").Eq(goqu.I(constants.QuizzesTable+".id")))).
		InnerJoin(goqu.T(constants.QuizQuestionsTable), goqu.On(goqu.I(constants.QuizzesTable+".id").Eq(goqu.I(constants.QuizQuestionsTable+".quiz_id")))).
		Where(goqu.Ex{
			UserPlayedQuizTable + ".user_id": userId,
			UserPlayedQuizTable + ".is_host": false,
		}).GroupBy("user_played_quizzes.id", "quizzes.id").Order(goqu.I("user_played_quizzes.created_at").Desc())

	if titleSearch != "" {
		query = query.Where(goqu.I("quizzes.title").ILike("%" + titleSearch + "%"))
	}

	offset := (page - 1) * constants.DefaultPageSize
	query = query.
		Limit(constants.DefaultPageSize).
		Offset(uint(offset))

	sql, args, err := query.ToSQL()
	if err != nil {
		return userPlayedQuiz, 0, err
	}

	err = model.db.ScanStructs(&userPlayedQuiz, sql, args...)
	if err != nil {
		return userPlayedQuiz, 0, err
	}

	var totalCount int64
	if len(userPlayedQuiz) > 0 {
		totalCount = userPlayedQuiz[0].TotalCount
	}

	return userPlayedQuiz, totalCount, err
}

func (model *UserPlayedQuizModel) ListUserPlayedQuizesWithQuestionById(UserPlayedQuizId string) ([]structs.ResUserPlayedQuizAnalyticsBoard, error) {
	var userPlayedQuizAnalyticsBoard []structs.ResUserPlayedQuizAnalyticsBoard

	query := model.db.From(UserPlayedQuizTable).
		Select(
			goqu.I(constants.UserQuizResponsesTable+".question_id").As("question_id"),
			goqu.I(constants.UserQuizResponsesTable+".answers").As("selected_answer"),
			goqu.I(constants.QuestionsTable+".answers").As("correct_answer"),
			"calculated_score",
			"is_attend",
			"response_time",
			"calculated_points",
			"question",
			"options",
			"question_media",
			"options_media",
			"resource",
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

	userPlayedQuizUUID, err := uuid.Parse(UserPlayedQuizId)
	if err != nil {
		return nil, err
	}
	responseModel := &UserQuizResponseModel{db: model.db}

	for index := 0; index < len(userPlayedQuizAnalyticsBoard); index++ {
		row := &userPlayedQuizAnalyticsBoard[index]
		options, correctAnswer, selectedAnswer, err := applyOptionOrderForRow(responseModel, userPlayedQuizUUID, row.QuestionID, row.Options, row.CorrectAnswer, row.SelectedAnswer)
		if err != nil {
			return nil, err
		}
		row.Options = options
		row.CorrectAnswer = correctAnswer
		row.SelectedAnswer = selectedAnswer
	}

	return userPlayedQuizAnalyticsBoard, nil
}

// JoinedUser model
type JoinedUser struct {
	UserID    string `json:"user_id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	ImageKey  string `json:"img_key" db:"img_key"`
}

func (model *UserPlayedQuizModel) GetJoinedUsers(activeQuizId string) ([]JoinedUser, error) {
	users := []JoinedUser{}

	err := model.db.Select(goqu.I("u.id"), goqu.L("COALESCE(upq.player_name, u.first_name)").As("first_name"), goqu.I("u.img_key")).
		From(goqu.T(UserPlayedQuizTable).As("upq")).
		Join(goqu.T(UserTable).As("u"), goqu.On(goqu.I("u.id").Eq(goqu.I("upq.user_id")))).
		Where(goqu.Ex{"upq.active_quiz_id": activeQuizId, "upq.is_host": false}).
		Order(goqu.I("upq.created_at").Asc()).
		ScanStructs(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (model *UserPlayedQuizModel) GetCountOfTotalJoinUsers(activeQuizId string) (int64, error) {
	return model.db.From(UserPlayedQuizTable).Where(goqu.Ex{
		"active_quiz_id": activeQuizId,
		"is_host":        false,
	}).Count()
}

// Deletes all user-played quizzes and associated responses for a specific user (userId)
func (model *UserPlayedQuizModel) DeleteUserPlayedQuizzesAndReponseByUserId(transaction *goqu.TxDatabase, userId string) error {

	subquery := transaction.From(UserPlayedQuizTable).Select("id").Where(goqu.Ex{"user_id": userId})

	_, err := transaction.Delete(UserQuizResponsesTable).Where(goqu.Ex{"user_played_quiz_id": goqu.Op{"in": subquery}}).Executor().Exec()
	if err != nil {
		return err
	}

	_, err = transaction.Delete(UserPlayedQuizTable).Where(goqu.Ex{"user_id": userId}).Executor().Exec()

	return err
}

// Get streak count using userPlayedQuizId and questionId
// use previous questionId to get previous streak count
func (model *UserPlayedQuizModel) GetStreakCount(userPlayedQuizId, questionId uuid.UUID) (int, error) {

	var previuosQuestionId string
	var streakCount int

	// Get previous questionId
	previuosQuestionFound, err := model.db.From(constants.QuizQuestionsTable).
		Select("question_id").
		Where(
			goqu.Ex{
				"next_question": questionId,
			},
		).ScanVal(&previuosQuestionId)
	if err != nil {
		return streakCount, err
	}

	// streak count is 0 if previous question is not found
	if !previuosQuestionFound {
		return streakCount, nil
	}

	found, err := model.db.From(UserQuizResponsesTable).
		Select("streak_count").
		Where(
			goqu.Ex{
				"user_played_quiz_id": userPlayedQuizId,
				"question_id":         previuosQuestionId,
			},
		).ScanVal(&streakCount)

	if err != nil {
		return streakCount, err
	}

	if !found {
		return streakCount, sql.ErrNoRows
	}

	return streakCount, nil
}
