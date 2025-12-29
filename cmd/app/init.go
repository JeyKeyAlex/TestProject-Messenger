package main

import (
	"github.com/JeyKeyAlex/TestProject-Messenger/internal/config"
	tpGRPC "github.com/JeyKeyAlex/TestProject-Messenger/internal/transport/grpc"
	tpGRPCMessenger "github.com/JeyKeyAlex/TestProject-Messenger/internal/transport/grpc/messenger"
	"github.com/rs/zerolog"
	googlegrpc "google.golang.org/grpc"
	"net"
	"time"
)

func initGRPCClientConnection(appConfig *config.Configuration) (*googlegrpc.ClientConn, error) {
	clientInfo := appConfig.ClientsGRPC.TestProject
	timeout := googlegrpc.WithIdleTimeout(clientInfo.IdleTimeout)

	dialOptions := []googlegrpc.DialOption{
		timeout,
	}

	conn, err := googlegrpc.NewClient(
		clientInfo.GetFullAddress(),
		dialOptions...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initKitGRPC(appConfig *config.Configuration, netLogger zerolog.Logger, listenErr chan error) (*googlegrpc.Server, net.Listener) {

	grpcUserServer := tpGRPCMessenger.NewServer()

	grpcServer := googlegrpc.NewServer(
		googlegrpc.MaxRecvMsgSize(appConfig.GRPC.MaxRequestBodySize),
		googlegrpc.MaxSendMsgSize(appConfig.GRPC.MaxRequestBodySize),
	)

	pbUser.RegisterUserServer(grpcServer, grpcUserServer)

	l, err := net.Listen(appConfig.GRPC.Network, appConfig.GRPC.Address)
	if err != nil {
		netLogger.Fatal().Err(err).Msg("failed to init net.Listen for grpc")
	} else {
		netLogger.Info().Msg("successful net.Listen for grpc init")
	}

	go tpGRPC.RunGRPCServer(grpcServer, l, netLogger, listenErr)
	time.Sleep(10 * time.Millisecond)
	return grpcServer, l
}
