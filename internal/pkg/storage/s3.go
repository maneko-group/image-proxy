package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type s3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type s3Storage struct {
	client s3API
	bucket string
}

func New(client s3API, bucket string) Getter {
	return &s3Storage{client: client, bucket: bucket}
}

func (s *s3Storage) Get(ctx context.Context, key string) (*Object, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: new(s.bucket),
		Key:    new(key),
	})
	if err != nil {
		return nil, mapError(key, err)
	}

	length := int64(-1)
	if out.ContentLength != nil {
		length = *out.ContentLength
	}

	return &Object{
		Body:          out.Body,
		ContentType:   aws.ToString(out.ContentType),
		ContentLength: length,
	}, nil
}

func mapError(key string, err error) error {
	if apiErr, ok := errors.AsType[smithy.APIError](err); ok {
		switch apiErr.ErrorCode() {
		case "NoSuchKey", "NoSuchBucket", "NotFound":
			return fmt.Errorf("%w: %q", ErrNotFound, key)
		}
	}

	return fmt.Errorf("storage: get %q: %w", key, err)
}
