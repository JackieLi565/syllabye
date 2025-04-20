package emailer

import (
	"context"
	"fmt"
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/JackieLi565/syllabye/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type NoReplyEmailer interface {
	SendWelcomeEmail(ctx context.Context, to string, name string) error
}

type sesNoReply struct {
	log    logger.Logger
	client *ses.Client
	from   *string
}

func NewSesNoReply(log logger.Logger, client *ses.Client) *sesNoReply {
	domain := os.Getenv(config.Domain)
	if domain == "" {
		panic("env var domain not defined")
	}

	return &sesNoReply{
		log:    log,
		client: client,
		from:   aws.String(fmt.Sprintf("noreply@%s", domain)),
	}
}

func (s *sesNoReply) SendWelcomeEmail(ctx context.Context, to string, name string) error {
	template := s.welcomeEmailTemplate(name)

	_, err := s.client.SendEmail(ctx, &ses.SendEmailInput{
		Source: s.from,
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		// TODO: Move to AWS email templates https://docs.aws.amazon.com/ses/latest/APIReference-V2/API_CreateEmailTemplate.html
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String("Welcome to Syllabye"),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(template),
				},
			},
		},
	})
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to send welcome email to %s", to), logger.Err(err))

		return util.ErrInternal
	}

	return nil
}

func (s *sesNoReply) welcomeEmailTemplate(name string) string {
	return fmt.Sprintf("Hello %s\n\n Welcome to Syllabye", name)
}
