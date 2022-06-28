package pkg

import (
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.AddConfigPath("$HOME/.product") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	viper.ReadInConfig()                  // Find and read the config file

	viper.SetEnvPrefix("product")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Default settings.

	// server
	viper.SetDefault("server.port", 80)

	// ddb
	viper.SetDefault("ddb.endpoint", "")
	viper.SetDefault("ddb.table.categories", "")
	viper.SetDefault("ddb.table.products", "")

	// DynamoDB table names passed via environment
	ddbTableProducts = viper.GetString("ddb.table.products")
	ddbTableCategories = viper.GetString("ddb.table.categories")

	// Allow DDB endpoint to be overridden to support amazon/dynamodb-local
	ddbEndpointOverride = viper.GetString("ddb.endpoint")

	viper.SetDefault("ddb.web.root_url", "")
	viper.SetDefault("ddb.image.root_url", "")

	webRootURL = viper.GetString("web.root_url")
	imageRootURL = viper.GetString("image.root_url")
}
