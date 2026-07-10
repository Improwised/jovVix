package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
)

type QuizCategoryController struct {
	quizCategoryModel *models.QuizCategoryModel
	appConfig         *config.AppConfig
	logger            *zap.Logger
}

func InitQuizCategoryController(db *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*QuizCategoryController, error) {
	return &QuizCategoryController{
		quizCategoryModel: models.InitQuizCategoryModel(db),
		appConfig:         appConfig,
		logger:            logger,
	}, nil
}

// isPublicQuizAdmin reports whether the requester is one of the configured
// public-quiz admin emails. Categories are a public-catalog feature, so the
// same allowlist that gates publishing quizzes gates managing categories.
func (ctrl *QuizCategoryController) isPublicQuizAdmin(c *fiber.Ctx) bool {
	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
	return ok && ctrl.appConfig.Quiz.IsPublicQuizAdmin(user.Email)
}

func (ctrl *QuizCategoryController) parseCategoryRequest(c *fiber.Ctx) (string, error) {
	var categoryReq structs.ReqQuizCategory
	err := json.Unmarshal(c.Body(), &categoryReq)
	if err != nil {
		return "", err
	}

	validate := validator.New()
	err = validate.Struct(categoryReq)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(categoryReq.Name), nil
}

// ListCategories lists every quiz category.
// Unauthenticated; the homepage uses this to group public quizzes by category.
func (ctrl *QuizCategoryController) ListCategories(c *fiber.Ctx) error {
	categories, err := ctrl.quizCategoryModel.ListCategories()
	if err != nil {
		ctrl.logger.Error("error occured while listing quiz categories", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrListCategories)
	}
	return utils.JSONSuccess(c, http.StatusOK, categories)
}

func (ctrl *QuizCategoryController) CreateCategory(c *fiber.Ctx) error {
	if !ctrl.isPublicQuizAdmin(c) {
		return utils.JSONFail(c, http.StatusForbidden, constants.ErrUnauthorized)
	}

	name, err := ctrl.parseCategoryRequest(c)
	if err != nil {
		ctrl.logger.Error("validate req error", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	exists, err := ctrl.quizCategoryModel.CategoryExistsByName(name, "")
	if err != nil {
		ctrl.logger.Error("error occured while checking category name", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrCreateCategory)
	}
	if exists {
		return utils.JSONFail(c, http.StatusConflict, constants.ErrCategoryAlreadyExists)
	}

	category, err := ctrl.quizCategoryModel.CreateCategory(name)
	if err != nil {
		ctrl.logger.Error("error occured while creating quiz category", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrCreateCategory)
	}

	return utils.JSONSuccess(c, http.StatusCreated, category)
}

func (ctrl *QuizCategoryController) UpdateCategory(c *fiber.Ctx) error {
	if !ctrl.isPublicQuizAdmin(c) {
		return utils.JSONFail(c, http.StatusForbidden, constants.ErrUnauthorized)
	}

	categoryId := c.Params(constants.CategoryId)

	name, err := ctrl.parseCategoryRequest(c)
	if err != nil {
		ctrl.logger.Error("validate req error", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	exists, err := ctrl.quizCategoryModel.CategoryExistsByName(name, categoryId)
	if err != nil {
		ctrl.logger.Error("error occured while checking category name", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUpdateCategory)
	}
	if exists {
		return utils.JSONFail(c, http.StatusConflict, constants.ErrCategoryAlreadyExists)
	}

	err = ctrl.quizCategoryModel.UpdateCategory(categoryId, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.ErrCategoryNotFound)
		}
		ctrl.logger.Error("error occured while updating quiz category", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrUpdateCategory)
	}

	return utils.JSONSuccess(c, http.StatusOK, "category updated successfully")
}

// DeleteCategory removes a category. Quizzes assigned to it are not deleted;
// they simply become uncategorized (ON DELETE SET NULL).
func (ctrl *QuizCategoryController) DeleteCategory(c *fiber.Ctx) error {
	if !ctrl.isPublicQuizAdmin(c) {
		return utils.JSONFail(c, http.StatusForbidden, constants.ErrUnauthorized)
	}

	categoryId := c.Params(constants.CategoryId)

	err := ctrl.quizCategoryModel.DeleteCategory(categoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusNotFound, constants.ErrCategoryNotFound)
		}
		ctrl.logger.Error("error occured while deleting quiz category", zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrDeleteCategory)
	}

	return utils.JSONSuccess(c, http.StatusOK, "category deleted successfully")
}
