package queue

import (
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewQueueClient() *sqs.Client {
	region := os.Getenv(config.AWS_REGION)
	accessKey := os.Getenv(config.AWS_ACCESS_KEY)
	secretKey := os.Getenv(config.AWS_SECRET_ACCESS_KEY)
	endpoint := os.Getenv(config.AWS_SQS_ENDPOINT)

	if region == "" || accessKey == "" || secretKey == "" || endpoint == "" {
		panic("missing AWS credentials")
	}

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))
	cfg := aws.Config{
		Region:       region,
		Credentials:  creds,
		BaseEndpoint: aws.String(endpoint),
	}

	return sqs.NewFromConfig(cfg)
}
