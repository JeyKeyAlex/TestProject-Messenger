package main

import (
	"github.com/JeyKeyAlex/TestProject-Messenger/internal/config"
	"github.com/JeyKeyAlex/TestProject-Messenger/pkg/logger"
)

func main() {
	//конфиги
	appConfig := config.MustLoad()

	//логгеры
	apiLogger, netLogger, coreLogger, loggerCloser := logger.ApiNetCoreCloserLoggers(appConfig)
	defer logger.CloseLogger(loggerCloser)

}
