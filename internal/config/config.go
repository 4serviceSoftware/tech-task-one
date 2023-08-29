package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	// Server
	ServerPort            string        `env:"SERVER_PORT,required"`
	ServerIdleTimeout     time.Duration `env:"SERVER_IDLE_TIMEOUT,default=30s"`
	ServerReadTimeout     time.Duration `env:"SERVER_READ_TIMEOUT,default=10s"`
	ServerWriteTimeout    time.Duration `env:"SERVER_WRITE_TIMEOUT,default=10s"`
	ServerShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT,default=30s"`
	ServerMaxUploadSize   int           `env:"SERVER_MAX_UPLOAD_SIZE,default=128"` // in megabytes

	// Postgres DB
	DbHost string `env:"DB_HOST,required"`
	DbPort string `env:"DB_PORT,required"`
	DbUser string `env:"DB_USER,required"`
	DbPass string `env:"DB_PASS,required"`
	DbName string `env:"DB_NAME,required"`

	// Static variables
	NodesCacheFilename string `env:"NODES_CACHE_FILENAME,default='./.cache/nodescache'"`
}

func GetFromEnv(ctx context.Context) (*Config, error) {
	config := Config{}
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
