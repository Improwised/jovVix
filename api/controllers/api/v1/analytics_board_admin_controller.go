package v1

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AnalyticsBoardAdminController struct {
	AnalyticsBoardAdminModel *models.AnalyticsBoardAdminModel
	presignedURLSvc          *services.PresignURLService
	logger                   *zap.Logger
	event                    *events.Events
}

func NewAnalyticsBoardAdminController(goqu *goqu.Database, logger *zap.Logger, event *events.Events, appConfig *config.AppConfig) (*AnalyticsBoardAdminController, error) {
	analyticsBoardAdminModel, err := models.InitAnalyticsBoardAdminModel(goqu)
	if err != nil {
		return nil, err
	}

	presignedURLSvc, err := services.NewFileUploadServices(appConfig.AWS.BucketName)
	if err != nil {
		return nil, err
	}

	return &AnalyticsBoardAdminController{
		AnalyticsBoardAdminModel: &analyticsBoardAdminModel,
		presignedURLSvc:          presignedURLSvc,
		logger:                   logger,
		event:                    event,
	}, nil

}

// GetAnalytics to send final score after quiz over to admin
// swagger:route GET /v1/analytics_board/admin AnalyticsBoardForAdmin RequestAnalyticsBoardForAdmin
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

	var wg sync.WaitGroup
	urlChan := make(chan URLResult, len(analyticsBoardData)*2)

	for i, v := range analyticsBoardData {
		if v.QuestionsMedia == "image" {
			wg.Add(1)
			go func(i int, resource string) {
				defer wg.Done()
				presignedURL, err := fc.presignedURLSvc.GetPresignedURL(resource, 5*time.Minute)
				if err != nil {
					fc.logger.Error("error while generating presign url for question media", zap.Error(err))
					urlChan <- URLResult{i, "", "", err}
					return
				}
				urlChan <- URLResult{i, "", presignedURL, nil}
			}(i, v.Resource)
		}

		if v.OptionsMedia == "image" {
			for optionKey, optionValue := range v.Options {
				wg.Add(1)
				go func(i int, optionKey string, optionValue string) {
					defer wg.Done()
					presignedURL, err := fc.presignedURLSvc.GetPresignedURL(optionValue, 1*time.Minute)
					if err != nil {
						fc.logger.Error("error while generating presign url for option media", zap.Error(err))
						urlChan <- URLResult{i, optionKey, "", err}
						return
					}
					urlChan <- URLResult{i, optionKey, presignedURL, nil}
				}(i, optionKey, optionValue)
			}
		}
	}

	go func() {
		wg.Wait()
		close(urlChan)
	}()

	for result := range urlChan {
		if result.err == nil && result.index < len(analyticsBoardData) {
			// For Question media
			if analyticsBoardData[result.index].QuestionsMedia == "image" {
				analyticsBoardData[result.index].Resource = result.url
			}
			// For Options media
			if analyticsBoardData[result.index].OptionsMedia == "image" {
				if result.optionKey != "" {
					analyticsBoardData[result.index].Options[result.optionKey] = result.url
				}
			}
		} else if result.err != nil {
			fc.logger.Error("Failed to update URL", zap.Error(result.err))
		}
	}

	return utils.JSONSuccess(ctx, http.StatusOK, analyticsBoardData)

}
