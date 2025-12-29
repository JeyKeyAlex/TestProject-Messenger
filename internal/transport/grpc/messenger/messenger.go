package messenger

import (
	"context"
	"github.com/JeyKeyAlex/TestProject-Messenger/internal/transport/grpc/common"
	pb "github.com/JeyKeyAlex/TestProject-genproto/messenger"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type RPCServer struct {
	create kitgrpc.Handler

	pb.UnimplementedMessengerServiceServer
}

func NewServer() pb.MessengerServiceServer {
	createEndpoint := makeCreate()
	return &RPCServer{
		create: kitgrpc.NewServer(createEndpoint, common.DecodeRequest, common.EncodeResponse),
	}
}

func makeCreate() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//reqID, ctx := middleware.GetRequestID(ctx)
		//serviceLogger := s.GetLogger().With().Str("func", "makeCreate").Str("request_id", reqID).Logger()
		//serviceLogger.Info().Msg("calling s.createUser")
		//
		//req, err := validate.CastValidateRequest[*pb.CreateUserRequest](s.GetValidator(), request)
		//if err != nil {
		//	serviceLogger.Error().Stack().Err(error_templates.ErrorDetailFromError(err)).Msg(pkgErr.FailedCastRequest)
		//	return nil, err
		//}

		return &pb.CreateMessageResponse{}, nil
	}
}
