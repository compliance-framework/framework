package runner

import (
	"context"
	"github.com/chris-cmsoft/concom/runner/proto"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.RunnerClient }

func (m *GRPCClient) Configure(config map[string]string) error {
	_, err := m.client.Configure(context.Background(), &proto.ConfigureRequest{Config: config})
	return err
}

func (m *GRPCClient) PrepareForEval() error {
	_, err := m.client.PrepareForEval(context.Background(), &proto.PrepareForEvalRequest{})
	return err
}

func (m *GRPCClient) Eval(bundlePath string) error {
	_, err := m.client.Eval(context.Background(), &proto.EvalRequest{BundlePath: bundlePath})
	return err
}

type GRPCServer struct {
	Impl Runner
}

func (m *GRPCServer) Configure(ctx context.Context, req *proto.ConfigureRequest) (*proto.ConfigureResponse, error) {
	return &proto.ConfigureResponse{}, m.Impl.Configure(req.Config)
}

func (m *GRPCServer) PrepareForEval(ctx context.Context, req *proto.PrepareForEvalRequest) (*proto.PrepareForEvalResponse, error) {
	return &proto.PrepareForEvalResponse{}, m.Impl.PrepareForEval()
}

func (m *GRPCServer) Eval(ctx context.Context, req *proto.EvalRequest) (*proto.EvalResponse, error) {
	return &proto.EvalResponse{}, m.Impl.Eval(req.BundlePath)
}
