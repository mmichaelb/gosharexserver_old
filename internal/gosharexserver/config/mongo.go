package config

import (
	"github.com/spf13/viper"
	"time"
)

func setMongoDefaults() {
	// connect process
	viper.SetDefault("mongodb.address", "localhost:27017")
	viper.SetDefault("mongodb.connect_timeout", time.Second*4)
	// authentication
	viper.SetDefault("mongodb.auth_db", "")
	viper.SetDefault("mongodb.auth_user", "")
	viper.SetDefault("mongodb.auth_passwd", "")
	// database
	viper.SetDefault("mongodb.db", "gosharexserver")
	// GridFS
	viper.SetDefault("mongodb.gridfs_prefix", "uploads")
	viper.SetDefault("mongodb.gridfs_chunk_size", 255000)
}
