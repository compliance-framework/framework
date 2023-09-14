package jobs

import (
	"testing"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func TestProcess(t *testing.T) {
	testCases := []struct {
		name              string
		assessmentResults AssessmentResults
		CreateFn          func(id string, object interface{}) error
	}{
		{
			name: "creates-result",
			assessmentResults: AssessmentResults{
				Id: "1234",
			},
			CreateFn: func(id string, object interface{}) error { return nil },
		}}

	f := FakeDriver{}
	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f.CreateFn = tc.CreateFn
			testCase := testCases[i]
			processJob := &ProcessJob{
				Driver: &f,
				ch:     make(chan *nats.Msg),
				Log:    zap.NewExample().Sugar(),
			}
			err := processJob.SaveAssessmentResults(testCase.assessmentResults)

			if err != nil {
				t.Errorf("failed to save assessment result: %s", err)
			}
		})
	}
}
