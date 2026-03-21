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

	return Config{
		Env:         env,
		Addr:        addr,
		DatabaseURL: databaseURL,
		S3Endpoint:  s3Endpoint,
	}, nil
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
