package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetAgentConfigFile returns the config file path from the command line arguments
//
// For simplicity and early feedback, this will panic if it can't find the config file as
// we need a config file to connect to NATS say, and it's pointless to run without this.
func GetAgentConfigFile(cmd *cobra.Command) string {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatal("Missing Config File")
	}
	if os.Stat(configFile); os.IsNotExist(err) {
		log.Fatal("Cannot find file")
	}
	return configFile
}

// ReadAgentConfig reads the agent config file and returns the config type. If it can't
// determine the config type, it will panic.
func ReadAgentConfig(configPath string) {
	fmt.Println(filepath.Base(configPath))
	viper.AddConfigPath(filepath.Dir(configPath))
	switch filepath.Ext(configPath) {
	case ".json":
		viper.SetConfigType("json")
		viper.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), ".json"))
	case ".yaml":
		viper.SetConfigType("yaml")
		viper.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), ".yaml"))
	case ".yml":
		viper.SetConfigType("yaml")
		viper.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), ".yml"))
	case ".toml":
		viper.SetConfigType("toml")
		viper.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), ".toml"))
	default:
		log.Fatalf("Unsupported config file type: %v", filepath.Ext(configPath))
	}

	viper.ReadInConfig()
}
