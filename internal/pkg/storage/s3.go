package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	*s3.S3
}

type Storage interface {
	Upload(ctx context.Context, key string, data []byte) (string, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
}

var (
	client *Client
	once   sync.Once
)

// DefaultClient return a default Client
func DefaultClient() Storage {
	if client == nil {
		once.Do(func() {
			cfg := config.GetConfig().Storage

			s3Config := &aws.Config{
				Endpoint:         &cfg.Endpoint,
				Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
				S3ForcePathStyle: aws.Bool(true),
				Region:           aws.String(cfg.Region),
			}
			s, err := session.NewSession(s3Config)
			if err != nil {
				logger.SugaredLogger.Panicw("Init S3 session failed", "err", err)
			}
			client = &Client{
				s3.New(s),
			}
		})
	}
	return client
}

// Upload file to S3 and return the proxy url and error
// the return url is the cloudflare proxy url, if err is not nil, the url is ""
func (cli *Client) Upload(ctx context.Context, key string, data []byte) (string, error) {
	cfg := config.GetConfig().Storage
	input := &s3.PutObjectInput{
		Bucket: aws.String(cfg.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	_, err := cli.PutObjectWithContext(ctx, input)
	// Return the proxy url, like "https://s3.domain.com/{bucketName}/{objectName}
	url := fmt.Sprintf("%s/%s/%s", cfg.Domain, cfg.BucketName, key)
	return url, err
}

func (cli *Client) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := cli.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(config.GetConfig().Storage.BucketName),
		Key:    aws.String(key),
	})
	return out.Body, err
}
