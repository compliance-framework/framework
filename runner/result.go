package runner

import "github.com/chris-cmsoft/concom/runner/proto"

type Result struct {
	Status       proto.ExecutionStatus   `json:"status"`
	AssessmentId string                  `json:"assessment_id"`
	Error        error                   `json:"error"`
	Observations []*proto.Observation    `json:"observations"`
	Findings     []*proto.Finding        `json:"findings"`
	Risks        []*proto.Risk           `json:"risks"`
	Logs         []*proto.LogEntry       `json:"logs"`
}

func ErrorResult(assessmentId string, err error) *Result {
	return &Result{
		Status:       proto.ExecutionStatus_FAILURE,
		AssessmentId: assessmentId,
		Error:        err,
		Observations: []*proto.Observation{},
		Findings:     []*proto.Finding{},
		Risks:        []*proto.Risk{},
		Logs:         []*proto.LogEntry{},
	}
}
