package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func loadViper(fileName string) (err error) {
	// set configuration filepath to the provided parameter
	viper.SetConfigFile(fileName)
	// add default values if the given config file does not contains specific values or do not exist
	// default values taken from /configs/default-config.toml
	// set webserver defaults
	setWebserverDefaults()
	// set MongoDB settings
	setMongoDefaults()
	// read config from filepath
	return viper.ReadInConfig()
}

func setWebserverDefaults() {
	viper.SetDefault("webserver.address", "localhost:10711")
	// reverse proxy header specifies whether a reverse proxy is used and the application should parse the remote ip
	viper.SetDefault("webserver.reverse_proxy_header", "")
	// whitelisted content types contains a list of all content types which should be displayed inline
	viper.SetDefault("webserver.whitelisted_content_types", []string{
		"image/png", "image/jpeg", "image/jpg", "image/gif",
		"text/plain", "text/plain; charset=utf-8",
		"video/mp4", "video/mpeg", "video/mpg4", "video/mpeg4", "video/flv",
	})
	// length of each generated authorization token is set to 20 per default
	viper.SetDefault("webserver.authorization_token_length", 20)
}

// LoadMainConfig loads the main config and stores the data into the Cfg variable.
func LoadMainConfig(fileName string) (err error) {
	if err = loadViper(fileName); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Could not read configuration from file, %T: %v. Falling back to defaults.\n", err, err)
			err = nil
		}
	}
	return
}
