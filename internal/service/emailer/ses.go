package emailer

import (
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

func NewSesClient() *ses.Client {
	region := os.Getenv(config.AWS_REGION)
	accessKey := os.Getenv(config.AWS_ACCESS_KEY)
	secretKey := os.Getenv(config.AWS_SECRET_ACCESS_KEY)

	if region == "" || accessKey == "" || secretKey == "" {
		panic("missing AWS credentials")
	}

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))
	cfg := aws.Config{
		Region:      region,
		Credentials: creds,
	}
	if os.Getenv(config.ENV) == "development" {
		endpoint := os.Getenv(config.AWS_SES_ENDPOINT)
		if endpoint == "" {
			panic("missing ses endpoint, required in development")
		}

		cfg.BaseEndpoint = aws.String(endpoint)
	}

	return ses.NewFromConfig(cfg)
}
