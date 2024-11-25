package internal

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetYamlPolicyPath(t *testing.T) {
	yaml := []byte(`nats:
  domain: localhost
  port: 4222

plugins:
  local-ssh-security:
    assessment_plan_id: "uuid"

    plugin: "ghcr.io/compliance-framework/plugin/local-ssh:0.0.5-upload"
    policies:
      - "../plugin-local-ssh-policies/dist/bundle.tar.gz"

    config:
      host: "machine1"
      username: "user"
      password: "password"

  local-ssh-security2:
    assessment_plan_id: "uuid2"

    schedule: "* * * * *"

    source: "ghcr.io/compliance-framework/plugin/local-ssh:0.0.5-upload"
    policies:
      - "../plugin-local-ssh-policies/dist/bundle.tar.gz"

    config:
      host: "machine2"
      username: "user"
      key: "/home/user/.ssh/id_rsa"

verbose: 2`)

	var agentConfig AgentConfig
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(yaml))

	err := viper.Unmarshal(&agentConfig)

	assert.Nil(t, err)
	assert.Equal(t, agentConfig.PluginConfig["local-ssh-security"].AssessmentPlanId, "uuid")
	assert.Equal(t, agentConfig.PluginConfig["local-ssh-security2"].AssessmentPlanId, "uuid2")
	assert.Equal(t, agentConfig.PluginConfig["local-ssh-security"].Schedule, "")
	assert.Equal(t, agentConfig.PluginConfig["local-ssh-security2"].Schedule, "* * * * *")
}
