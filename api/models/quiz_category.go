package models

import (
	"database/sql"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

const QuizCategoriesTable = "quiz_categories"

// QuizCategoryModel implements quiz category related database operations
type QuizCategoryModel struct {
	db *goqu.Database
}

// InitQuizCategoryModel initializes the QuizCategoryModel
func InitQuizCategoryModel(goquDB *goqu.Database) *QuizCategoryModel {
	return &QuizCategoryModel{db: goquDB}
}

type QuizCategory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ListCategories returns every quiz category ordered by name.
func (model *QuizCategoryModel) ListCategories() ([]QuizCategory, error) {
	categories := []QuizCategory{}
	err := model.db.From(QuizCategoriesTable).
		Select("id", "name", "created_at", "updated_at").
		Order(goqu.I("name").Asc()).
		ScanStructs(&categories)
	return categories, err
}

func (model *QuizCategoryModel) GetCategoryById(categoryId string) (QuizCategory, error) {
	var category QuizCategory
	found, err := model.db.From(QuizCategoriesTable).
		Select("id", "name", "created_at", "updated_at").
		Where(goqu.Ex{"id": categoryId}).
		Limit(1).
		ScanStruct(&category)
	if err != nil {
		return category, err
	}
	if !found {
		return category, sql.ErrNoRows
	}
	return category, nil
}

// CategoryExistsByName reports whether a category with the given name already
// exists (case-insensitive). excludeId skips one category, useful for renames.
func (model *QuizCategoryModel) CategoryExistsByName(name string, excludeId string) (bool, error) {
	query := model.db.From(QuizCategoriesTable).
		Where(goqu.Func("LOWER", goqu.C("name")).Eq(strings.ToLower(strings.TrimSpace(name))))

	if excludeId != "" {
		query = query.Where(goqu.C("id").Neq(excludeId))
	}

	count, err := query.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (model *QuizCategoryModel) CreateCategory(name string) (QuizCategory, error) {
	var category QuizCategory

	categoryId, err := uuid.NewUUID()
	if err != nil {
		return category, err
	}

	found, err := model.db.Insert(QuizCategoriesTable).Rows(
		goqu.Record{
			"id":   categoryId,
			"name": name,
		},
	).Returning("id", "name", "created_at", "updated_at").Executor().ScanStruct(&category)
	if err != nil {
		return category, err
	}
	if !found {
		return category, sql.ErrNoRows
	}

	return category, nil
}

func (model *QuizCategoryModel) UpdateCategory(categoryId, name string) error {
	result, err := model.db.Update(QuizCategoriesTable).
		Set(goqu.Record{
			"name":       name,
			"updated_at": goqu.L("now()"),
		}).
		Where(goqu.Ex{"id": categoryId}).
		Executor().Exec()
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteCategory removes a category; quizzes referencing it fall back to
// uncategorized via the ON DELETE SET NULL foreign key.
func (model *QuizCategoryModel) DeleteCategory(categoryId string) error {
	result, err := model.db.Delete(QuizCategoriesTable).
		Where(goqu.Ex{"id": categoryId}).
		Executor().Exec()
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
