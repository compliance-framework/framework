package runner

import (
	"context"
	"github.com/chris-cmsoft/concom/runner/proto"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.RunnerClient }

func (m *GRPCClient) PrepareForEval() error {
	_, err := m.client.PrepareForEval(context.Background(), &proto.PrepareForEvalRequest{})
	return err
}

type GRPCServer struct {
	Impl Runner
}

func (m *GRPCServer) PrepareForEval(ctx context.Context, req *proto.PrepareForEvalRequest) (*proto.PrepareForEvalResponse, error) {
	return &proto.PrepareForEvalResponse{}, m.Impl.PrepareForEval()
}
