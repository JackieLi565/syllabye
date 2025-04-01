package bucket

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/service/logger"
)

func NewS3Client(log logger.Logger) *s3.Client {
	region := os.Getenv(config.AWS_REGION)
	accessKey := os.Getenv(config.AWS_ACCESS_KEY)
	secretKey := os.Getenv(config.AWS_SECRET_ACCESS_KEY)
	endpoint := os.Getenv(config.AWS_S3_ENDPOINT)

	if region == "" || accessKey == "" || secretKey == "" || endpoint == "" {
		log.Error("missing AWS credentials")
	}

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))

	cfg := aws.Config{
		Region:       region,
		Credentials:  creds,
		BaseEndpoint: aws.String(endpoint),
	}

	return s3.NewFromConfig(cfg)
}
