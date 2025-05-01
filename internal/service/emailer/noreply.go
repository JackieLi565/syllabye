package emailer

import (
	"context"
	"encoding/json"
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
	SendSubmissionSuccessEmail(ctx context.Context, to string, name string, course string) error
	SendSubmissionMissingEmail(ctx context.Context, to string, name string, course string) error
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
	welcomeTemplate := os.Getenv(config.AWS_SES_WELCOME_TEMPLATE)
	if welcomeTemplate == "" {
		s.log.Error("Welcome template name not defined")
		return util.ErrInternal
	}

	templateData := map[string]interface{}{
		"name": name,
	}

	return s.sendEmail(ctx, to, welcomeTemplate, templateData)
}

func (s *sesNoReply) SendSubmissionSuccessEmail(ctx context.Context, to string, name string, course string) error {
	uploadSuccessTemplate := os.Getenv(config.AWS_SES_UPLOAD_SUCCESS_TEMPLATE)
	if uploadSuccessTemplate == "" {
		s.log.Error("Upload Success template name not defined")
		return util.ErrInternal
	}

	templateData := map[string]interface{}{
		"name":   name,
		"course": course,
	}

	return s.sendEmail(ctx, to, uploadSuccessTemplate, templateData)
}

func (s *sesNoReply) SendSubmissionMissingEmail(ctx context.Context, to string, name string, course string) error {
	uploadErrorTemplate := os.Getenv(config.AWS_SES_UPLOAD_ERROR_TEMPLATE)
	if uploadErrorTemplate == "" {
		s.log.Error("Upload Error template name not defined")
		return util.ErrInternal
	}

	templateData := map[string]interface{}{
		"name":   name,
		"course": course,
		"reason": "We did not receive your upload file",
	}

	return s.sendEmail(ctx, to, uploadErrorTemplate, templateData)
}

func (s *sesNoReply) sendEmail(ctx context.Context, to string, template string, templateData map[string]interface{}) error {
	dat, _ := json.Marshal(templateData)

	res, err := s.client.SendTemplatedEmail(ctx, &ses.SendTemplatedEmailInput{
		Source: s.from,
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Template:     aws.String(template),
		TemplateData: aws.String(string(dat)),
	})
	if err != nil {
		s.log.Error(fmt.Sprintf("failed to send welcome email to %s", to), logger.Err(err))
		return util.ErrInternal
	}

	s.log.Info(fmt.Sprintf("email template %s sent to %s with message id %s", template, to, *res.MessageId))

	return err
}
