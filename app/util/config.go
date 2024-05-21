package util

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Config *viper.Viper

// InitConfig initializes the configuration from file and environment variables.
func InitConfig(configPath string) error {

	envPrefix := os.Getenv("ENV_PREFIX")
	envConfigPath, isExists := os.LookupEnv("ENV_CONFIG_PATH")
	if isExists {
		configPath = envConfigPath
	}

	Config = viper.New()
	Config.SetConfigType("yaml")
	Config.SetConfigFile(configPath)
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	Config.SetEnvPrefix(envPrefix)

	err := Config.ReadInConfig()
	return err
}
