package main

import (
	"net"
	"time"

	"github.com/JeyKeyAlex/TestProject-Messenger/internal/config"
	tpGRPC "github.com/JeyKeyAlex/TestProject-Messenger/internal/transport/grpc"
	tpGRPCMessenger "github.com/JeyKeyAlex/TestProject-Messenger/internal/transport/grpc/messenger"

	"github.com/rs/zerolog"
	googlegrpc "google.golang.org/grpc"

	pb "github.com/JeyKeyAlex/TestProject-genproto/messenger"
)

func initKitGRPC(appConfig *config.Configuration, netLogger zerolog.Logger, listenErr chan error) (*googlegrpc.Server, net.Listener) {

	grpcUserServer := tpGRPCMessenger.NewServer()

	grpcServer := googlegrpc.NewServer(
		googlegrpc.MaxRecvMsgSize(appConfig.GRPC.MaxRequestBodySize),
		googlegrpc.MaxSendMsgSize(appConfig.GRPC.MaxRequestBodySize),
	)

	pb.RegisterMessengerServiceServer(grpcServer, grpcUserServer)

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
