package services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PresignURLService struct {
	s3Client *s3.Client
	bucket   string
}

func NewFileUploadServices(bucket string) (*PresignURLService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &PresignURLService{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func (fuSvc *PresignURLService) GetPresignedURL(objectKey string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(fuSvc.s3Client)

	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(fuSvc.bucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
