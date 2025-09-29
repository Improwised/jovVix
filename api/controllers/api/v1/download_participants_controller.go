package v1

import (
	"net/http"
	"strings"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DownloadParticipantsController struct {
	QuizModel                 *models.QuizModel
	FinalScoreBoardAdminModel models.FinalScoreBoardAdminModel
	Config                    *config.AppConfig
}

func InitDownloadParticipantsController(db *goqu.Database, logger *zap.Logger, appConfig *config.AppConfig) (*DownloadParticipantsController, error) {
	quizModel := models.InitQuizModel(db)

	finalScoreBoardAdminModel, err := models.InitFinalScoreBoardAdminModel(db)
	if err != nil {
		return nil, err
	}

	return &DownloadParticipantsController{
		QuizModel:                 quizModel,
		FinalScoreBoardAdminModel: finalScoreBoardAdminModel,
		Config:                    appConfig,
	}, nil
}

func (dcp *DownloadParticipantsController) DownloadParticipantsReport(c *fiber.Ctx) error {
	activeQuizId := c.Params(constants.ActiveQuizId)
	isTop10Users := strings.ToLower(c.Query(constants.IsTop10Users))
	userOrder := strings.ToLower(c.Query(constants.SortOrder))

	order := constants.DescOrder

	quizAnalysis, err := dcp.QuizModel.GetQuizAnalysis(activeQuizId)
	if err != nil {
		return utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	totalUsers, err := dcp.FinalScoreBoardAdminModel.TotoalUsers(activeQuizId)
	if err != nil {
		return err
	}

	limit := totalUsers

	if isTop10Users == constants.Top10UserTrueValue {
		limit = 10
	}

	if userOrder == constants.AscOrder {
		order = constants.AscOrder
	}

	rankData, err := dcp.FinalScoreBoardAdminModel.GetScoreForAdmin(activeQuizId, limit, 0, order)
	if err != nil {
		return err
	}

	input := utils.InputParticipants{
		QuizAnalysis:         quizAnalysis,
		FinalScoreBoardAdmin: rankData,
	}

	filePath, err := services.GenerateCsv(dcp.Config.FilePath, activeQuizId, input, func(ip utils.InputParticipants) [][]string {
		return utils.GetParticipantsData(ip)
	})

	if err != nil {
		return err
	}
	return utils.CsvFileResponse(c, filePath, activeQuizId+".csv")
}
