package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"go.uber.org/zap"
)

type AnalyticsBoardAdminController struct {
	AnalyticsBoardAdminModel *models.AnalyticsBoardAdminModel
	presignedURLSvc          *services.PresignURLService
	logger                   *zap.Logger
	appConfig                *config.AppConfig
}

func NewAnalyticsBoardAdminController(goqu *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*AnalyticsBoardAdminController, error) {
	analyticsBoardAdminModel, err := models.InitAnalyticsBoardAdminModel(goqu)
	if err != nil {
		return nil, err
	}

	presignedURLSvc, err := services.NewFileUploadServices(&appConfig.AWS)
	if err != nil {
		return nil, err
	}

	return &AnalyticsBoardAdminController{
		AnalyticsBoardAdminModel: &analyticsBoardAdminModel,
		presignedURLSvc:          presignedURLSvc,
		logger:                   logger,
		appConfig:                appConfig,
	}, nil

}

// GetAnalytics to send final score after quiz over to admin
// swagger:route GET /v1/analytics_board/admin AnalyticsBoard RequestAnalyticsBoardForAdmin
//
// Get a analyticsboard details for admin.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseAnalyticsBoardForAdmin
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (fc *AnalyticsBoardAdminController) GetAnalyticsForAdmin(ctx *fiber.Ctx) error {

	var activeQuizId = ctx.Query(constants.ActiveQuizId, "")

	if activeQuizId == "" || !(len(activeQuizId) == 36) {
		fc.logger.Error("active quiz id is not valid")
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("active quiz id should be valid string").Error())
	}

	analyticsBoardData, err := fc.AnalyticsBoardAdminModel.GetAnalyticsForAdmin(activeQuizId)
	if err != nil {
		fc.logger.Error("Error while getting analytics for admin", zap.Error(err))
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}

	services.ProcessAnalyticsData(analyticsBoardData, fc.presignedURLSvc, fc.logger)

	return utils.JSONSuccess(ctx, http.StatusOK, analyticsBoardData)

}

func (fc *AnalyticsBoardAdminController) DownloadQuizReport(c *fiber.Ctx) error {
	activeQuizId := c.Params(constants.ActiveQuizId)
	if activeQuizId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrQuizIdNotFound)
	}

	quizReport, err := fc.AnalyticsBoardAdminModel.GetAnalyticsForAdmin(activeQuizId)

	if err != nil {
		return err
	}

	orderToUserAndQuestionData := make(map[int][]structs.UserAndQuestionData)

	for _, data := range quizReport {
		user := structs.UserAndQuestionData{
			Options:        data.Options,
			UserName:       data.UserName,
			CorrectAnswer:  data.CorrectAnswer,
			SelectedAnswer: data.SelectedAnswer.String,
			Question:       data.Question,
			QuestionType:   data.QuestionType,
		}
		if val, ok := orderToUserAndQuestionData[data.OrderNo]; ok {
			orderToUserAndQuestionData[data.OrderNo] = append(val, user)
		} else {
			orderToUserAndQuestionData[data.OrderNo] = []structs.UserAndQuestionData{user}
		}
	}

	orderOfQuestion := make([]int, 0, len(orderToUserAndQuestionData))
	for o := range orderToUserAndQuestionData {
		orderOfQuestion = append(orderOfQuestion, o)
	}
	sort.Ints(orderOfQuestion)

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)
	utils.BuildHeading(m)

	utils.BuildUsersTables(m, orderToUserAndQuestionData, orderOfQuestion)
	pdfPath := filepath.Join(fc.appConfig.PDFS_FILE_PATH, fmt.Sprintf("/%s.pdf", activeQuizId))
	err = m.OutputFileAndClose(pdfPath)

	if err != nil {
		return err
	}

	data, err := os.ReadFile(pdfPath)
	if err != nil {
		log.Printf("Error reading PDF file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading PDF file")
	}
	return utils.JSONSuccessPdf(c, http.StatusOK, data)
}
