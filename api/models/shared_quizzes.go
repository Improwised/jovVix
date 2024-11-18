package models

import (
	"database/sql"
	"net/url"

	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const SharedQuizzesTable = "shared_quizzes"

type SharedQuizzes struct {
	Id         string `json:"id" db:"id"`
	QuizId     string `json:"quiz_id" db:"quiz_id"`
	SharedTo   string `json:"shared_to" db:"shared_to"`
	SharedBy   string `json:"shared_by" db:"shared_by"`
	Permission string `json:"permission" db:"permission"`
}

type SharedQuizzesModel struct {
	db     *goqu.Database
	logger *zap.Logger
}

func InitSharedQuizzesModel(goquDB *goqu.Database, logger *zap.Logger) *SharedQuizzesModel {
	return &SharedQuizzesModel{db: goquDB, logger: logger}
}

// Insert the data for share quiz
func (model *SharedQuizzesModel) InsertSharedQuiz(sharedQuizzes SharedQuizzes) error {
	Id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	_, err = model.db.Insert(SharedQuizzesTable).Rows(
		goqu.Record{
			"id":         Id,
			"quiz_id":    sharedQuizzes.QuizId,
			"shared_to":  sharedQuizzes.SharedTo,
			"shared_by":  sharedQuizzes.SharedBy,
			"permission": sharedQuizzes.Permission,
		},
	).Executor().Exec()
	if err != nil {
		return err
	}

	return nil
}

// List authorized users for perticular quiz
func (model *SharedQuizzesModel) ListQuizAuthorizedUsersByQuizId(quizId string) ([]structs.ResUserWithQuizPermission, error) {
	var usersWithAccess []structs.ResUserWithQuizPermission

	err := model.db.From(SharedQuizzesTable).
		Select(
			goqu.I("shared_quizzes.id"),
			goqu.I("shared_quizzes.shared_to"),
			goqu.I("users.first_name"),
			goqu.I("users.last_name"),
			goqu.I("users.img_key"),
			goqu.I("shared_quizzes.permission"),
		).
		LeftJoin(goqu.T("users"), goqu.On(goqu.I("shared_quizzes.shared_to").Eq(goqu.I("users.email")))).
		Where(goqu.Ex{"quiz_id": quizId}).
		ScanStructs(&usersWithAccess)

	return usersWithAccess, err
}

// Update authorized user permission for perticular quiz
func (model *SharedQuizzesModel) UpdateUserPermissionById(id string, reqUpdatePermission structs.ReqShareQuiz) error {

	_, err := model.db.Update(SharedQuizzesTable).Set(goqu.Record{
		"permission": reqUpdatePermission.Permission,
		"updated_at": goqu.L("now()"),
	}).Where(goqu.Ex{"id": id}).Executor().Exec()

	return err
}

// Delete authorized user permission for perticular quiz
func (model *SharedQuizzesModel) DeleteUserPermissionById(id string) error {

	_, err := model.db.Delete(SharedQuizzesTable).Where(goqu.Ex{"id": id}).Executor().Exec()

	return err
}

// List shared quiz for perticular user (only shared with the user or shared by the user)
func (model *SharedQuizzesModel) ListSharedQuizzes(sharedBy, sharedTo string) ([]QuizWithQuestions, error) {

	questionsCountSubquery := model.db.From("quiz_questions").
		Select(goqu.COUNT(goqu.DISTINCT("question_id"))).
		Where(goqu.C("quiz_id").Eq(goqu.I("quizzes.id")))

	query := model.db.From("quizzes").
		Select(
			goqu.I("quizzes.id"),
			goqu.I("quizzes.title"),
			goqu.I("quizzes.description"),
			goqu.I("quizzes.creator_id"),
			goqu.I("quizzes.created_at"),
			goqu.I("quizzes.updated_at"),
			questionsCountSubquery.As("total_questions"),
		).
		InnerJoin(
			goqu.T("shared_quizzes"),
			goqu.On(goqu.I("quizzes.id").Eq(goqu.I("shared_quizzes.quiz_id"))),
		).
		LeftJoin(
			goqu.T("quiz_questions"),
			goqu.On(goqu.I("quiz_questions.quiz_id").Eq(goqu.I("quizzes.id"))),
		)

	// Apply filters only if the parameters are provided
	if sharedTo != "" {
		query = query.Where(goqu.I("shared_quizzes.shared_to").Eq(sharedTo))
	}
	if sharedBy != "" {
		query = query.Where(goqu.I("shared_quizzes.shared_by").Eq(sharedBy))
	}

	// Add grouping and ordering
	query = query.GroupBy(
		goqu.I("quizzes.id"),
		goqu.I("quizzes.title"),
		goqu.I("quizzes.description"),
		goqu.I("quizzes.creator_id"),
		goqu.I("quizzes.created_at"),
		goqu.I("quizzes.updated_at"),
	).Order(goqu.I("quizzes.created_at").Desc())

	rows, err := query.Executor().Query()
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

// Check the user is creator of the perticular quiz or not
func (model *SharedQuizzesModel) CheckQuizCreatorExists(quizId, creatorId string) (bool, error) {
	found, err := model.db.From("quizzes").
		Where(
			goqu.Ex{"id": quizId},
			goqu.Ex{"creator_id": creatorId},
		).
		Select(goqu.L("1")).Executor().ScanVal(new(int))

	if err != nil {
		return false, err
	}

	return found, nil
}

// Get permission of quiz for perticular user
func (model *SharedQuizzesModel) GetPermissionByQuizAndUser(quizID, sharedTo string) (string, error) {
	var permission string

	found, err := model.db.From("shared_quizzes").
		Select("permission").
		Where(
			goqu.Ex{"quiz_id": quizID},
			goqu.Ex{"shared_to": sharedTo},
		).
		Executor().
		ScanVal(&permission)

	if err != nil {
		return permission, err
	}

	if !found {
		return permission, sql.ErrNoRows
	}

	return permission, nil
}
