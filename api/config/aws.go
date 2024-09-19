package config

type AWSConfig struct {
	BucketName       string `envconfig:"BUCKET_NAME"`
	Region           string `envconfig:"AWS_REGION"`
	S3BucketEndpoint string `envconfig:"S3_BUCKET_ENDPOINT"`
}
