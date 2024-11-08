package runner

import (
	"github.com/chris-cmsoft/concom/runner/proto"
	"github.com/google/uuid"
	"testing"
)

func TestCallableEvalResponse_AddFinding(t *testing.T) {
	resp := NewCallableEvalResponse()

	if len(resp.Findings) > 0 {
		t.Errorf("len(resp.Findings): got %d, want %d", len(resp.Findings), 0)
	}

	findingId := uuid.New().String()
	resp.AddFinding(&proto.Finding{
		Id:    findingId,
		Title: "A rather brilliant finding",
	})

	if len(resp.Findings) != 1 {
		t.Errorf("len(resp.Findings): got %d, want %d", len(resp.Findings), 1)
	}

	if resp.Findings[0].Id != findingId {
		t.Errorf("resp.Findings[0].Id: got %s, want %s", resp.Findings[0].Id, findingId)
	}
}

func TestCallableEvalResponse_AddObservation(t *testing.T) {
	resp := NewCallableEvalResponse()

	if len(resp.Observations) > 0 {
		t.Errorf("len(resp.Findings): got %d, want %d", len(resp.Observations), 0)
	}

	observationId := uuid.New().String()
	resp.AddObservation(&proto.Observation{
		Id:    observationId,
		Title: "Some clever observation",
	})

	if len(resp.Observations) != 1 {
		t.Errorf("len(resp.Findings): got %d, want %d", len(resp.Observations), 1)
	}

	if resp.Observations[0].Id != observationId {
		t.Errorf("resp.Findings[0].Id: got %s, want %s", resp.Observations[0].Id, observationId)
	}
}

func TestCallableEvalResponse_Result(t *testing.T) {
	resp := NewCallableEvalResponse()

	if resp.Result() != resp.EvalResponse {
		t.Errorf("resp.Result(): got %v, want %v", resp.Result(), resp.EvalResponse)
	}
}
