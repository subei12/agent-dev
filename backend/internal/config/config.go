package config

import (
	"fmt"
	"os"
)

type Config struct {
	Env         string
	Addr        string
	DatabaseURL string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
}

func Load() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	s3Endpoint := os.Getenv("S3_ENDPOINT")
	if s3Endpoint == "" {
		s3Endpoint = "http://localhost:9000"
	}

	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	if s3AccessKey == "" {
		s3AccessKey = "minio"
	}

	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	if s3SecretKey == "" {
		s3SecretKey = "minio123"
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		s3Bucket = "agent-platform-dev"
	}

	return Config{
		Env:         env,
		Addr:        addr,
		DatabaseURL: databaseURL,
		S3Endpoint:  s3Endpoint,
		S3AccessKey: s3AccessKey,
		S3SecretKey: s3SecretKey,
		S3Bucket:    s3Bucket,
	}, nil
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
