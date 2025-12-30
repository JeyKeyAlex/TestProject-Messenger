package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/JeyKeyAlex/TestProject-Messenger/internal/config"
	"github.com/JeyKeyAlex/TestProject-Messenger/pkg/logger"
)

func main() {
	appConfig := config.MustLoad()

	coreLogger, loggerCloser := logger.CoreCloserLoggers(appConfig)
	defer logger.CloseLogger(loggerCloser)

	listenErr := make(chan error, 1)

	grpcServer, grpcListener := initKitGRPC(appConfig, coreLogger, listenErr)
	defer func() {
		err := grpcListener.Close()
		if err != nil {
			coreLogger.Fatal().Err(err).Msg("error closing grpc listener")
		}
	}()

	var shutdownCh = make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var err error

	for {
		select {
		case err = <-listenErr:
			if err != nil {
				coreLogger.Error().Err(err).Msg("received listener error")
				shutdownCh <- os.Kill
			}
		case sig := <-shutdownCh:
			coreLogger.Info().Msgf("received shutdown signal: %v", sig)
			grpcServer.GracefulStop()
			coreLogger.Info().Msg("server loop stopped")
			return
		}
	}
}
