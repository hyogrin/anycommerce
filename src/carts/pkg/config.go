package pkg

import (
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")       // name of config file (without extension)
	viper.AddConfigPath("$HOME/.carts") // call multiple times to add many search paths
	viper.AddConfigPath(".")            // optionally look for config in the working directory
	viper.ReadInConfig()                // Find and read the config file

	viper.SetEnvPrefix("carts")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Default settings.

	// server
	viper.SetDefault("server.port", 80)
}
