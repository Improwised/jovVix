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

type AnalyticsBoardUserController struct {
	analyticsBoardUserModel *models.AnalyticsBoardUserModel
	presignedURLSvc         *services.PresignURLService
	logger                  *zap.Logger
	event                   *events.Events
}

type URLResult struct {
	index     int
	optionKey string
	url       string
	err       error
}

func NewAnalyticsBoardUserController(goqu *goqu.Database, logger *zap.Logger, event *events.Events, appConfig *config.AppConfig) (*AnalyticsBoardUserController, error) {
	analyticsBoardUserModel, err := models.InitAnalyticsBoardUserModel(goqu)
	if err != nil {
		return nil, err
	}

	presignedURLSvc, err := services.NewFileUploadServices(appConfig.AWS.BucketName)
	if err != nil {
		return nil, err
	}

	return &AnalyticsBoardUserController{
		analyticsBoardUserModel: &analyticsBoardUserModel,
		presignedURLSvc:         presignedURLSvc,
		logger:                  logger,
		event:                   event,
	}, nil

}

// GetAnalytics to send analytics board after quiz over to user
// swagger:route GET /v1/analytics_board/user AnalyticsBoardForUser RequestAnalyticsBoardForUser
//
// Get a analyticsboard details for user.
//
//		Consumes:
//		- application/json
//
//		Schemes: http, https
//
//		Responses:
//		  200: ResponseAnalyticsBoardForUser
//	     400: GenericResFailNotFound
//		  500: GenericResError
func (fc *AnalyticsBoardUserController) GetAnalyticsForUser(ctx *fiber.Ctx) error {
	userPlayedQuizId := ctx.Query(constants.UserPlayedQuiz)
	fc.logger.Debug("AnalyticsBoardUserController.GetAnalyticsForUser called", zap.Any("userPlayedQuizId", userPlayedQuizId))

	if userPlayedQuizId == "" || !(len(userPlayedQuizId) == 36) {
		fc.logger.Error("user played quiz is not valid")
		return utils.JSONFail(ctx, http.StatusBadRequest, errors.New("user play quiz should be valid string").Error())
	}

	fc.logger.Debug("analyticsBoardUserModel.GetAnalyticsForUser called", zap.Any("userPlayedQuizId", userPlayedQuizId))
	analyticsBoardData, err := fc.analyticsBoardUserModel.GetAnalyticsForUser(userPlayedQuizId)
	if err != nil {
		fc.logger.Error("Error while getting analytics for user", zap.Error(err))
		return utils.JSONFail(ctx, http.StatusInternalServerError, errors.New("internal server error").Error())
	}
	fc.logger.Debug("analyticsBoardUserModel.GetAnalyticsForUser success", zap.Any("analyticsBoardData", analyticsBoardData))

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

	fc.logger.Debug("AnalyticsBoardUserController.GetAnalyticsForUser success", zap.Any("analyticsBoardData", analyticsBoardData))

	return utils.JSONSuccess(ctx, http.StatusOK, analyticsBoardData)
}
