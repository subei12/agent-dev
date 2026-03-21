package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStore struct {
	client *minio.Client
	bucket string
}

func NewMinIO(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return minio.New(parsed.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: parsed.Scheme == "https" || strings.HasPrefix(endpoint, "https://"),
	})
}

func NewObjectStore(endpoint, accessKey, secretKey, bucket string) (*ObjectStore, error) {
	client, err := NewMinIO(endpoint, accessKey, secretKey)
	if err != nil {
		return nil, err
	}

	return &ObjectStore{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *ObjectStore) PutJSON(ctx context.Context, objectKey string, value any) error {
	if err := s.ensureBucket(ctx); err != nil {
		return err
	}

	body, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(ctx, s.bucket, objectKey, bytes.NewReader(body), int64(len(body)), minio.PutObjectOptions{
		ContentType: "application/json",
	})
	return err
}

func (s *ObjectStore) PutBytes(ctx context.Context, objectKey string, data []byte, contentType string) error {
	if err := s.ensureBucket(ctx); err != nil {
		return err
	}

	_, err := s.client.PutObject(ctx, s.bucket, objectKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (s *ObjectStore) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
}
