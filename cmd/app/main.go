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

	//testProjectConn, err := initGRPCClientConnection(appConfig)
	//if err != nil {
	//	panic(err)
	//}

	listenErr := make(chan error, 1)

	grpcServer, grpcListener := initKitGRPC(appConfig, netLogger, listenErr)

}
