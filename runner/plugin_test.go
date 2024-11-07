package runner

import "testing"

func TestRunnerConfig_GetString(t *testing.T) {
	config := &RunnerConfig{
		ConfigItem{
			Key:   "host",
			Value: "192.168.1.1",
		},
		ConfigItem{
			Key:   "port",
			Value: 80,
		},
	}
	if config.GetString("host") != "192.168.1.1" {
		t.Errorf("Expected to be able to read configuration")
	}
	if config.GetInt("port") != 80 {
		t.Errorf("Expected to be able to read configuration")
	}
}
