package services

import (
	"io"
	"path/filepath"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/gofiber/storage/s3"
	"go.uber.org/zap"
)

type FileUploadService struct {
	logger  *zap.Logger
	event   *events.Events
	pub     *watermill.WatermillPublisher
	config  *config.AWSConfig
	storage *s3.Storage
}

func NewFileUploadService(logger *zap.Logger, event *events.Events, pub *watermill.WatermillPublisher, config *config.AWSConfig) *FileUploadService {

	storage := s3.New(s3.Config{
		Bucket:   config.BucketName,
		Region:   config.Region,
		Endpoint: config.S3BucketEndpoint,
	})

	return &FileUploadService{
		logger:  logger,
		event:   event,
		pub:     pub,
		config:  config,
		storage: storage,
	}
}

// UploadToS3 uploads a file to the S3 bucket
func (fuSvc *FileUploadService) UploadToS3(folder string, filename string, fileData io.Reader) error {

	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		return err
	}

	filePath := filepath.Join(folder, filename)

	return fuSvc.storage.Set(filePath, fileBytes, 1*time.Minute)
}
