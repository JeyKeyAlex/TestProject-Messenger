package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func MustLoad() *Configuration {
	var envFiles []string
	if _, err := os.Stat(".env"); err == nil {
		log.Println("found .env file, adding it to env config files list")
		envFiles = append(envFiles, ".env")
	}
	if os.Getenv("APP_ENV") != "" {
		appEnvName := fmt.Sprintf(".env.%s", os.Getenv("APP_ENV"))
		if _, err := os.Stat(appEnvName); err == nil {
			log.Println("found", appEnvName, "file, adding it to env config files list")
			envFiles = append(envFiles, appEnvName)
		}
	}
	if len(envFiles) > 0 {
		err := godotenv.Overload(envFiles...)
		if err != nil {
			log.Fatalf("error while opening env config: %s", err.Error())
		}
	}
	cfg := &Configuration{}
	ctx := context.Background()

	err := envconfig.Process(ctx, cfg)
	if err != nil {
		log.Fatalf("error while parsing env config: %s", err.Error())
	}
	return cfg
}

type (
	Configuration struct {
		Log         LogConfig   `env:",prefix=LOG_"`
		ClientsGRPC ClientsGRPC `env:",prefix=CLIENTS_GRPC_"`
		Version     Version     `env:",prefix=VERSION_"`
		GRPC        GRPCConfig  `env:",prefix=GRPC_"`
	}

	LogConfig struct {
		Level             string        `env:"LEVEL,default=info"`
		Batch             bool          `env:"BATCH,default=false"`
		BatchSize         int           `env:"BATCH_SIZE,default=1000"`
		BatchPollInterval time.Duration `env:"BATCH_POLL_INTERVAL,default=5s"`
	}
	Version struct {
		Number string `env:"NUMBER,default=1.0.0"`
		Build  string `env:"BUILD,default=dev"`
	}

	GRPCConfig struct {
		RequestLoggingEnabled      bool          `env:"REQUEST_LOGGING_ENABLED,default=false"`
		ResponseTimeLoggingEnabled bool          `env:"RESPONSE_TIME_LOGGING_ENABLED,default=false"`
		ReadTimeout                time.Duration `env:"READ_TIMEOUT,default=30s"`
		WriteTimeout               time.Duration `env:"WRITE_TIMEOUT,default=30s"`
		IdleTimeout                time.Duration `env:"IDLE_TIMEOUT,default=30s"`
		MaxRequestBodySize         int           `env:"MAX_REQUEST_BODY_SIZE,default=33554432"`
		Network                    string        `env:"NETWORK,default=tcp"`
		Address                    string        `env:"ADDRESS,default=:18080"`
	}

	ClientsGRPC struct {
		TestProject ClientGRPC `env:",prefix=TEST_PROJECT"`
	}

	ClientGRPC struct {
		Address             string        `env:"ADDRESS"`
		Port                string        `env:"PORT"`
		IdleTimeout         time.Duration `env:"IDLE_TIMEOUT"`
		InsecureSkipVerify  bool          `env:"INSECURE_SKIP_VERIFY,default=false"`
		MaxRequestBodySize  int           `env:"MAX_REQUEST_BODY_SIZE,default=4194304"`
		MaxResponseBodySize int           `env:"MAX_RESPONSE_BODY_SIZE,default=4194304"`
	}
)

func (c ClientGRPC) GetFullAddress() string {
	return c.Address + c.Port
}
