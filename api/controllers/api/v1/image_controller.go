package v1

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/events"
	"github.com/Improwised/jovvix/api/pkg/watermill"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ImageController struct {
	questionModel *models.QuestionModel
	fileUploadSvc *services.FileUploadService
	logger        *zap.Logger
	event         *events.Events
	pub           *watermill.WatermillPublisher
	config        *config.AppConfig
}

func NewImageController(goqu *goqu.Database, logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher, config *config.AppConfig) (*ImageController, error) {

	questionModel := models.InitQuestionModel(goqu, logger)
	fileUploadSvc := services.NewFileUploadService(logger, event, pub, &config.AWS)

	return &ImageController{
		questionModel: questionModel,
		fileUploadSvc: fileUploadSvc,
		logger:        logger,
		event:         event,
		pub:           pub,
		config:        config,
	}, nil
}

// InsertImage to insert image in s3 bucket
// swagger:route POST /v1/images FileUpload RequestInsertImage
//
// Insert image in s3 bucket.
//
//			Consumes:
//			- multipart/form-data
//
//			Schemes: http, https
//
//			Responses:
//			  200: ResponseInsertImage
//		     400: GenericResFailNotFound
//	     401: GenericResFailConflict
//			  500: GenericResError
func (imgc *ImageController) InsertImage(c *fiber.Ctx) error {

	quizId := c.Query("quiz_id")
	if quizId == "" {
		return utils.JSONError(c, http.StatusBadRequest, "No quiz_id found")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return utils.JSONError(c, http.StatusBadRequest, "No files found")
	}

	// Get the files from the form
	files := form.File["image-attachment"]
	if len(files) == 0 {
		return utils.JSONError(c, http.StatusBadRequest, "No files found")
	}

	errChan := make(chan error, len(files))
	var wg sync.WaitGroup

	for _, fileHeader := range files {
		wg.Add(1)
		go func(fileHeader *multipart.FileHeader) {
			defer wg.Done()

			file, err := fileHeader.Open()
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()

			// Upload to S3
			err = imgc.fileUploadSvc.UploadToS3(quizId, fileHeader.Filename, file)
			if err != nil {
				errChan <- err
				return
			}

			splitStringArr := strings.Split(fileHeader.Filename, "_")
			if len(splitStringArr) == 1 {
				err = imgc.questionModel.UpdateResourceOfQuestionById(fileHeader.Filename, filepath.Join(quizId, fileHeader.Filename))
				if err != nil {
					errChan <- err
					return
				}
			}

			if len(splitStringArr) == 2 {
				err = imgc.questionModel.UpdateOptionsOfQuestionById(splitStringArr[1], splitStringArr[0], filepath.Join(quizId, fileHeader.Filename))
				if err != nil {
					errChan <- err
					return
				}
			}
		}(fileHeader)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			imgc.logger.Error(constants.ErrInsertImage, zap.Error(err))
			return utils.JSONError(c, http.StatusInternalServerError, constants.ErrInsertImage)
		}
	}

	return utils.JSONSuccess(c, http.StatusOK, "Images uploaded successfully!")
}
