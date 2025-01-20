package runner

import (
	"context"
	proto2 "github.com/compliance-framework/framework/runner/proto"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto2.RunnerClient }

func (m *GRPCClient) Configure(req *proto2.ConfigureRequest) (*proto2.ConfigureResponse, error) {
	return m.client.Configure(context.Background(), req)
}

func (m *GRPCClient) PrepareForEval(req *proto2.PrepareForEvalRequest) (*proto2.PrepareForEvalResponse, error) {
	return m.client.PrepareForEval(context.Background(), req)
}

func (m *GRPCClient) Eval(req *proto2.EvalRequest) (*proto2.EvalResponse, error) {
	resp, err := m.client.Eval(context.Background(), req)
	return resp, err
}

type GRPCServer struct {
	Impl Runner
}

func (m *GRPCServer) Configure(ctx context.Context, req *proto2.ConfigureRequest) (*proto2.ConfigureResponse, error) {
	return m.Impl.Configure(req)
}

func (m *GRPCServer) PrepareForEval(ctx context.Context, req *proto2.PrepareForEvalRequest) (*proto2.PrepareForEvalResponse, error) {
	return m.Impl.PrepareForEval(req)
}

func (m *GRPCServer) Eval(ctx context.Context, req *proto2.EvalRequest) (*proto2.EvalResponse, error) {
	return m.Impl.Eval(req)
}
