package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/JackieLi565/syllabye/internal/service/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type PresignerClient interface {
	GetObject(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error)
	PutObject(ctx context.Context, objectKey string, contentType string, checksum string, lifetimeSecs int64) (string, error)
}

type s3Presigner struct {
	presignClient *s3.PresignClient
	log           logger.Logger
	bucket        string
}

func NewS3Presigner(log logger.Logger, s3Client *s3.Client, bucket string) *s3Presigner {
	return &s3Presigner{
		log:           log,
		presignClient: s3.NewPresignClient(s3Client),
		bucket:        bucket,
	}
}

func (p *s3Presigner) GetObject(ctx context.Context, objectKey string, lifetimeSecs int64) (string, error) {
	request, err := p.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		p.log.Error(fmt.Sprintf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			p.bucket, objectKey, err))
	}

	return request.URL, err
}

func (p *s3Presigner) PutObject(ctx context.Context, objectKey string, contentType string, checksum string, lifetimeSecs int64) (string, error) {
	request, err := p.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:            aws.String(p.bucket),
		Key:               aws.String(objectKey),
		ContentType:       aws.String(contentType),
		ChecksumAlgorithm: types.ChecksumAlgorithmCrc32,
		ChecksumCRC32:     aws.String(checksum),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		p.log.Error(fmt.Sprintf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			p.bucket, objectKey, err))
	}

	return request.URL, err
}
