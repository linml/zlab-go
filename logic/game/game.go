package game

import (
	"context"
	hello "zlab/protobuf"
)

type HelloService struct {
}

func (h *HelloService) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	resp := new(hello.HelloResponse)

	return resp, nil
}
