package common

import (
	"context"
)

func DecodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}
func EncodeResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}
