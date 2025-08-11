package v1

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"go.uber.org/zap"
)

type DownloadQuizReportController struct {
	AnalyticsBoardAdminModel *models.AnalyticsBoardAdminModel
	db                       *goqu.Database
	logger                   *zap.Logger
}

func NewDownloadQuizReportContorller(goqu *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*DownloadQuizReportController, error) {
	analyticsBoardAdminModel, err := models.InitAnalyticsBoardAdminModel(goqu)
	if err != nil {
		return nil, err
	}

	return &DownloadQuizReportController{
		AnalyticsBoardAdminModel: &analyticsBoardAdminModel,
		db:                       goqu,
		logger:                   logger,
	}, nil
}

func (dqrc *DownloadQuizReportController) DownloadQuizReport(c *fiber.Ctx) error {
	activeQuizId := c.Params(constants.ActiveQuizId)
	if activeQuizId == "" {
		return utils.JSONError(c, http.StatusBadRequest, constants.ErrKratosIDEmpty)
	}

	quizReport, err := dqrc.AnalyticsBoardAdminModel.GetAnalyticsForAdmin(activeQuizId)

	if err != nil {
		return err
	}

	questionsAndUsersMap := make(map[int][]structs.Users)
	var questions []string

	for _, data := range quizReport {
		user := structs.Users{
			Options:        data.Options,
			UserName:       data.UserName,
			CorrectAnswer:  data.CorrectAnswer,
			SelectedAnswer: data.SelectedAnswer.String,
		}
		if val, ok := questionsAndUsersMap[data.OrderNo]; ok {
			questionsAndUsersMap[data.OrderNo] = append(val, user)
		} else {
			questionsAndUsersMap[data.OrderNo] = []structs.Users{user}
			questions = append(questions, data.Question)
		}
	}

	keys := make([]int, 0, len(questionsAndUsersMap))
	for k := range questionsAndUsersMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)
	utils.BuildHeading(m)

	utils.BuildUsersTables(m, questionsAndUsersMap, questions, keys)
	pdfPath := fmt.Sprintf("%s/%s.pdf", config.GetConfigByName("PDFS_FILE_PATH"), activeQuizId)
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
