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
		Log                     LogConfig                     `env:",prefix=LOG_"`
		Runtime                 RuntimeConfig                 `env:",prefix=RUNTIME_"`
		HealthCheck             HealthCheckConfig             `env:",prefix=HEALTHCHECK_"`
		Debug                   DebugConfig                   `env:",prefix=DEBUG_"`
		ClientsGRPC             ClientsGRPC                   `env:",prefix=CLIENTS_GRPC_"`
		HTTP                    HTTPConfig                    `env:",prefix=HTTP_"`
		UpcomingPaymentReminder UpcomingPaymentReminderConfig `env:",prefix=UPCOMING_PAYMENT_REMINDER_"`
		Version                 Version                       `env:",prefix=VERSION_"`
		Metrics                 MetricsConfig                 `env:",prefix=METRICS_"`
		OTel                    OTelConfig                    `env:",prefix=OPEN_TELEMETRY_"`
	}

	LogConfig struct {
		Level             string        `env:"LEVEL,default=info"`
		Batch             bool          `env:"BATCH,default=false"`
		BatchSize         int           `env:"BATCH_SIZE,default=1000"`
		BatchPollInterval time.Duration `env:"BATCH_POLL_INTERVAL,default=5s"`
	}

	RuntimeConfig struct {
		UseCPUs    int `env:"USE_CPUS,default=0"`
		MaxThreads int `env:"MAX_THREADS,default=0"`
	}

	HealthCheckConfig struct {
		GoroutineThreshold int `env:"GOROUTINE_THRESHOLD,default=20"`
	}

	DebugConfig struct {
	}
	Version struct {
		Number string `env:"NUMBER,default=1.0.0"`
		Build  string `env:"BUILD,default=dev"`
	}

	ClientsGRPC struct {
		RhumbAPI ClientGRPC `env:",prefix=RHUMB_API_"`
	}

	ClientGRPC struct {
		Address             string        `env:"ADDRESS"`
		Port                string        `env:"PORT"`
		IdleTimeout         time.Duration `env:"IDLE_TIMEOUT"`
		InsecureSkipVerify  bool          `env:"INSECURE_SKIP_VERIFY,default=false"`
		MaxRequestBodySize  int           `env:"MAX_REQUEST_BODY_SIZE,default=4194304"`
		MaxResponseBodySize int           `env:"MAX_RESPONSE_BODY_SIZE,default=4194304"`
	}

	HTTPConfig struct {
		CORSEnabled                bool          `env:"CORS_ENABLED,default=false"`
		RequestLoggingEnabled      bool          `env:"REQUEST_LOGGING_ENABLED,default=false"`
		ResponseTimeLoggingEnabled bool          `env:"RESPONSE_TIME_LOGGING_ENABLED,default=false"`
		ReadTimeout                time.Duration `env:"READ_TIMEOUT,default=30s"`
		WriteTimeout               time.Duration `env:"WRITE_TIMEOUT,default=30s"`
		IdleTimeout                time.Duration `env:"IDLE_TIMEOUT,default=30s"`
		MaxRequestBodySize         int           `env:"MAX_REQUEST_BODY_SIZE,default=4194304"`
		Network                    string        `env:"NETWORK,default=tcp"`
		Address                    string        `env:"ADDRESS,default=:8081"`
	}

	UpcomingPaymentReminderConfig struct {
		TickerInterval   time.Duration `env:"TICKER_INTERVAL,default=1h"`
		ReminderInterval time.Duration `env:"INTERVAL,default=24h"`
	}

	MetricsConfig struct {
		Path      string `env:"PATH,default=/metrics"`
		Namespace string `env:"NAMESPACE,default=rhumb"`
		Subsystem string `env:"SUBSYSTEM,default=notify-worker"`
	}
	OTelConfig struct {
		TracerName  string         `env:"TRACER_NAME,default=rhumb-notify-worker-tracer"`
		ServiceName string         `env:"SERVICE_NAME,default=rhumb-notify-worker"`
		Exporter    ExporterConfig `env:",prefix=EXPORTER_"`
		Sampler     SamplerConfig  `env:",prefix=SAMPLER_"`
	}

	ExporterConfig struct {
		URI                string        `env:"URI,default=localhost:4318"`
		URLPath            string        `env:"URL_PATH,default=/v1/traces"` // HTTP-only usage; for gRPC use localhost:4317 and .proto
		BatchTimeout       time.Duration `env:"BATCH_TIMEOUT,default=1s"`
		MaxExportBatchSize int32         `env:"MAX_EXPORT_BATCH_SIZE,default=512"`
		MaxQueueSize       int32         `env:"MAX_QUEUE_SIZE,default=2048"`
		WithInsecure       bool          `env:"WITH_INSECURE,default=true"`
	}

	SamplerConfig struct {
		Fraction float64 `env:"FRACTION,default=0.5"`
	}
)

func (c ClientGRPC) GetFullAddress() string {
	return c.Address + c.Port
}
