package services

import (
	"context"
	"reflect"
	"sync"
	"time"

	aws_config "github.com/Improwised/jovvix/api/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
)

type URLResult struct {
	index     int
	optionKey string
	url       string
	err       error
}

type PresignURLService struct {
	s3Client *s3.Client
	bucket   string
}

func NewFileUploadServices(awsConfig *aws_config.AWSConfig) (*PresignURLService, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:               awsConfig.S3BucketEndpoint,
				SigningRegion:     awsConfig.Region,
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsConfig.AwsAccessKeyId, awsConfig.AWSSecretAccessKey, "")),
		config.WithRegion(awsConfig.Region),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &PresignURLService{
		s3Client: s3Client,
		bucket:   awsConfig.BucketName,
	}, nil
}

// Generate presigned url using s3 client
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

// Processes a slice of analytics data items, fetching presigned URLs for media fields (if "QuestionsMedia" or "OptionsMedia" are image)
// and updating the respective fields in the data (update image_id to presigned url for display the image in UI)
func ProcessAnalyticsData[T any](data []T, presignedURLSvc *PresignURLService, logger *zap.Logger) {
	var wg sync.WaitGroup
	urlChan := make(chan URLResult, len(data)*2)

	for i, v := range data {
		vValue := reflect.ValueOf(v)

		if vValue.FieldByName("QuestionsMedia").String() == "image" {
			wg.Add(1)
			go func(i int, resource string) {
				defer wg.Done()
				presignedURL, err := presignedURLSvc.GetPresignedURL(resource, 5*time.Minute)
				if err != nil {
					logger.Error("error while generating presign url for question media", zap.Error(err))
					urlChan <- URLResult{i, "", "", err}
					return
				}
				urlChan <- URLResult{i, "", presignedURL, nil}
			}(i, vValue.FieldByName("Resource").String())
		}

		if vValue.FieldByName("OptionsMedia").String() == "image" {
			options := vValue.FieldByName("Options")
			for _, optionKey := range options.MapKeys() {
				wg.Add(1)
				go func(i int, optionKey string, optionValue string) {
					defer wg.Done()
					presignedURL, err := presignedURLSvc.GetPresignedURL(optionValue, 1*time.Minute)
					if err != nil {
						logger.Error("error while generating presign url for option media", zap.Error(err))
						urlChan <- URLResult{i, optionKey, "", err}
						return
					}
					urlChan <- URLResult{i, optionKey, presignedURL, nil}
				}(i, optionKey.String(), options.MapIndex(optionKey).String())
			}
		}
	}

	go func() {
		wg.Wait()
		close(urlChan)
	}()

	for result := range urlChan {
		if result.err == nil && result.index < len(data) {
			vValue := reflect.ValueOf(&data[result.index]).Elem()
			if vValue.FieldByName("QuestionsMedia").String() == "image" {
				vValue.FieldByName("Resource").SetString(result.url)
			}
			if vValue.FieldByName("OptionsMedia").String() == "image" && result.optionKey != "" {
				options := vValue.FieldByName("Options")
				options.SetMapIndex(reflect.ValueOf(result.optionKey), reflect.ValueOf(result.url))
			}
		} else if result.err != nil {
			logger.Error("Failed to update URL", zap.Error(result.err))
		}
	}
}
