package storage

import (
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

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
