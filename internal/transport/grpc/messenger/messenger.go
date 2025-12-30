package messenger

import (
	"context"
	"errors"
	"fmt"

	"github.com/JeyKeyAlex/TestProject-Messenger/internal/transport/grpc/common"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	pb "github.com/JeyKeyAlex/TestProject-genproto/messenger"
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
		req, ok := request.(*pb.CreateMessageRequest)
		if !ok {
			err := errors.New("invalid request fields")
			return nil, err
		}

		resp := fmt.Sprintf("%s, you are observed", req.Email)

		return &pb.CreateMessageResponse{
			Message: resp,
		}, nil
	}
}
