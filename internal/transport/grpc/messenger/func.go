package messenger

import (
	"context"

	pb "github.com/JeyKeyAlex/TestProject-genproto/messenger"
)

func (s *RPCServer) Create(ctx context.Context, req *pb.CreateMessageRequest) (*pb.CreateMessageResponse, error) {
	_, resp, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateMessageResponse), nil
}
