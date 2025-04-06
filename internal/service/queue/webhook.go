package queue

import (
	"context"
	"encoding/json"
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type WebhookMessage struct {
	RequestId string            `json:"requestId"`
	Url       string            `json:"url"`
	Method    string            `json:"method"`
	Payload   string            `json:"payload,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

type WebhookQueue interface {
	SendMessage(ctx context.Context, message WebhookMessage, delay int32) error
}

type sqsWebhook struct {
	log    logger.Logger
	client *sqs.Client
	url    string
}

func NewSqsWebhook(log logger.Logger, client *sqs.Client) *sqsWebhook {
	queueUrl := os.Getenv(config.AWS_SQS_WEBHOOK_URL)
	if queueUrl == "" {
		panic("missing webhook queue url")
	}

	return &sqsWebhook{
		log:    log,
		client: client,
		url:    queueUrl,
	}
}

func (s *sqsWebhook) SendMessage(ctx context.Context, message WebhookMessage, delay int32) error {
	messageBody, err := json.Marshal(message)
	if err != nil {
		return util.ErrMalformed
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		DelaySeconds: delay,
		MessageBody:  aws.String(string(messageBody)),
		QueueUrl:     aws.String(s.url),
	})
	if err != nil {
		s.log.Error("failed to create webhook queue message")
		return err
	}

	return nil
}
