package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/JeyKeyAlex/TestProject-Messenger/internal/config"
)

const (
	DefaultTimestampFieldName   = "time"
	DefaultLevelFieldName       = "level"
	DefaultMessageFieldName     = "message"
	DefaultErrorStackFieldName  = "stacktrace"
	DefaultCallerSkipFrameCount = 2
	componentField              = "component"
	Debug                       = "debug"
	debugMessage                = "debug log is enabled"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = DefaultTimestampFieldName
	zerolog.LevelFieldName = DefaultLevelFieldName
	zerolog.MessageFieldName = DefaultMessageFieldName
	zerolog.ErrorStackFieldName = DefaultErrorStackFieldName
	zerolog.CallerSkipFrameCount = DefaultCallerSkipFrameCount
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

type myWriteCloser struct {
	io.Writer
}

func (mwc *myWriteCloser) Close() error {
	log.Warn().Msg("Closing log writer")
	// Noop
	return nil
}

func ApiNetCoreCloserLoggers(cfg *config.Configuration) (zerolog.Logger, zerolog.Logger, zerolog.Logger, io.WriteCloser) {
	baseLog, logCloser := baseLogger(cfg)

	apiLogger := newComponentLogger(baseLog, "api")
	netLogger := newComponentLogger(baseLog, "net")
	coreLogger := newComponentLogger(baseLog, "core")

	return apiLogger, netLogger, coreLogger, logCloser
}

// CloseLogger не возвращает ошибку потому что она не обрабатывается и функция вызывается при остановке программы
func CloseLogger(loggerCloser io.WriteCloser) {
	if loggerCloser != nil {
		err := loggerCloser.Close()
		if err != nil {
			log.Error().Msg("error acquired while closing log writer: " + err.Error())
		}
	}
}

func baseLogger(cfg *config.Configuration) (zerolog.Logger, io.WriteCloser) {
	var err error
	var logger zerolog.Logger
	var loggerCloser io.WriteCloser
	if cfg.Log.Batch {
		logger, loggerCloser, err = newDiodeLogger(os.Stdout, cfg.Log.Level, cfg.Log.BatchSize, cfg.Log.BatchPollInterval)
	} else {
		logger, loggerCloser, err = newLogger(os.Stdout, cfg.Log.Level)
	}
	if err != nil {
		panic(err)
	}
	logger = logger.With().
		Str("worker", "rhumb-api").
		CallerWithSkipFrameCount(DefaultCallerSkipFrameCount).
		Logger()

	zerolog.TimestampFieldName = DefaultTimestampFieldName
	return logger, loggerCloser
}

func newLogger(w io.Writer, logLevel string) (zerolog.Logger, io.WriteCloser, error) {
	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return zerolog.Logger{}, nil, err
	}

	wc, ok := w.(io.WriteCloser)
	if !ok {
		wc = &myWriteCloser{w}
	}

	logger := zerolog.New(wc).Level(lvl).With().Timestamp().Logger()
	if lvl == zerolog.DebugLevel {
		logger.Debug().Bool(Debug, true).Msg(debugMessage)
	}

	return logger, wc, nil
}

func newDiodeLogger(w io.Writer, logLevel string, batchSize int, batchPollInterval time.Duration) (zerolog.Logger, io.WriteCloser, error) {
	logger, writer, err := newLogger(w, logLevel)
	if err != nil {
		return logger, writer, err
	}
	logWriter := diode.NewWriter(writer, batchSize, batchPollInterval, func(missed int) {
		logger.Warn().Str(componentField, "drop-logger").Int("dropped_count", missed).Msg("catch some dropped logs!")
	})
	logger = logger.Output(logWriter)
	return logger, logWriter, nil
}

func newComponentLogger(logger zerolog.Logger, component string) zerolog.Logger {
	return logger.With().Str(componentField, component).Logger()
}
